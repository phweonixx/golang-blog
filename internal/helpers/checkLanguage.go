package helpers

import "errors"

func CheckLanguage(lang string) error {
	validLanguages := map[string]bool{"en": true, "uk": true}
	if !validLanguages[lang] {
		return errors.New("invalid language: valid values are 'en' or 'uk'")
	}
	return nil
}
