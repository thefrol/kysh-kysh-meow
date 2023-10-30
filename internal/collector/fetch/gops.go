package fetch

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

const CPUpollInterval = 10 * time.Millisecond // время опроса процессора

func GoPS() Batcher {
	m, err := mem.VirtualMemory()
	if err != nil {
		// можем только вывести в журнал ошибку Mentor
		log.Error().Err(err).Msg("Ошибка чтения gopsutil.mem")
		return EmptyBatch{}
	}

	perCPU := true
	cpu, err := cpu.Percent(CPUpollInterval, perCPU) // todo я не уверен, что он тут выдает
	if err != nil {
		// можем только вывести в журнал ошибку Mentor
		log.Error().Err(err).Msg("Ошибка чтения gopsutil.cpu")
		return EmptyBatch{}
	}

	return PSBatch{
		vmStat:  m,
		cpuUtil: cpu,
	}
}

type PSBatch struct {
	vmStat  *mem.VirtualMemoryStat
	cpuUtil []float64
}

func (b PSBatch) ToTransport() (m []metrica.Metrica) {
	m = append(m, metrica.Gauge(b.vmStat.Total).Metrica("TotalMemory"))
	m = append(m, metrica.Gauge(b.vmStat.Free).Metrica("FreeMemory"))
	for i, cpu := range b.cpuUtil {
		label := fmt.Sprintf("%s%d", "CPUutilization", i)
		m = append(m, metrica.Gauge(cpu).Metrica(label))
	}

	return
}

type EmptyBatch struct{}

func (b EmptyBatch) ToTransport() []metrica.Metrica {
	return nil
}
