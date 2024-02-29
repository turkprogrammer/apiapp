//Код предоставляет структуру Validator для управления ошибками валидации.

package validator

// Validator представляет собой структуру для управления ошибками валидации.
type Validator struct {
	Errors      []string          `json:",omitempty"` // Общие ошибки валидации.
	FieldErrors map[string]string `json:",omitempty"` // Ошибки валидации для конкретных полей.
}

// HasErrors возвращает true, если есть какие-либо ошибки валидации.
func (v Validator) HasErrors() bool {
	return len(v.Errors) != 0 || len(v.FieldErrors) != 0
}

// AddError добавляет общую ошибку валидации.
func (v *Validator) AddError(message string) {
	if v.Errors == nil {
		v.Errors = []string{}
	}

	v.Errors = append(v.Errors, message)
}

// AddFieldError добавляет ошибку валидации для конкретного поля.
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}

	// Если для поля уже есть ошибка, она не будет перезаписана.
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Check добавляет общую ошибку валидации, если условие (ok) ложно.
func (v *Validator) Check(ok bool, message string) {
	if !ok {
		v.AddError(message)
	}
}

// CheckField добавляет ошибку валидации для конкретного поля, если условие (ok) ложно.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}
