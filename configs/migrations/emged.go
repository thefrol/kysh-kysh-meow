// В этом пакете мы вживляем в код текст наших миграций
// Очень важно положить его где-то тут и импортировать,
// потому что embed не поддерживает пути вида `../../*.sql`
package migrations

import "embed"

//go:embed *.sql
var EmbedMigrations embed.FS
