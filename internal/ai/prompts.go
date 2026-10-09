package ai

// ParticipantDetailsSystemPrompt ограничивает извлечение разделом участника проекта.
const ParticipantDetailsSystemPrompt = `Ты — система извлечения реквизитов из русскоязычных документов.
Извлекай данные ТОЛЬКО из раздела «УЧАСТНИК ПРОЕКТА» (включая содержимое этого раздела до начала следующего самостоятельного раздела).
Полностью игнорируй раздел «ОПЕРАТОР КЖНО» и любые другие разделы. Не используй значения из игнорируемых разделов даже если нужное поле там заполнено.
Не додумывай и не исправляй значения. Сохраняй написание реквизитов как в документе, удаляя только лишние пробелы по краям.
Все значения в fields должны быть строками. Если поле не найдено в разрешённом разделе, верни пустую строку.
summary — краткое описание найденных сведений об участнике проекта на русском языке. Если раздел не найден, укажи это в summary.
Верни только JSON, соответствующий переданной JSON Schema.`

// ParticipantDetailsSchema используется Ollama для ограничения структуры ответа.
var ParticipantDetailsSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"document_type": map[string]any{"type": "string"},
		"summary":       map[string]any{"type": "string"},
		"fields": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"organization":      map[string]any{"type": "string"},
				"legal_address":     map[string]any{"type": "string"},
				"ogrn":              map[string]any{"type": "string"},
				"inn":               map[string]any{"type": "string"},
				"kpp":               map[string]any{"type": "string"},
				"bank_account":      map[string]any{"type": "string"},
				"director_position": map[string]any{"type": "string"},
				"director_name":     map[string]any{"type": "string"},
			},
			"required":             []string{"organization", "legal_address", "ogrn", "inn", "kpp", "bank_account", "director_position", "director_name"},
			"additionalProperties": false,
		},
	},
	"required":             []string{"document_type", "summary", "fields"},
	"additionalProperties": false,
}
