package city_service

import (
	"regexp"
	"strings"
)

var (
	slugRegex = regexp.MustCompile(`[^a-z0-9]+`)
)

func (s *CityService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.TrimSpace(slug)

	// Транслитерация кириллицы (базовая)
	slug = transliterate(slug)

	// Замена всех не-латинских символов и пробелов на дефисы
	slug = slugRegex.ReplaceAllString(slug, "-")

	// Убираем дефисы по краям
	slug = strings.Trim(slug, "-")

	return slug
}

func transliterate(s string) string {
	translitMap := map[string]string{
		"а": "a", "б": "b", "в": "v", "г": "g", "д": "d",
		"е": "e", "ё": "yo", "ж": "zh", "з": "z", "и": "i",
		"й": "y", "к": "k", "л": "l", "м": "m", "н": "n",
		"о": "o", "п": "p", "р": "r", "с": "s", "т": "t",
		"у": "u", "ф": "f", "х": "kh", "ц": "ts", "ч": "ch",
		"ш": "sh", "щ": "shch", "ъ": "", "ы": "y", "ь": "",
		"э": "e", "ю": "yu", "я": "ya",
	}

	result := strings.Builder{}
	for _, r := range s {
		if replacement, ok := translitMap[string(r)]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
