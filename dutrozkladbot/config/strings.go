// Package config ..
package config

import "fmt"

// NotAvailableMessage ..
const NotAvailableMessage = "*Помилка*: на даний момент, розклад недоступний."

// ErrorMessage ..
const ErrorMessage = "e-rozklad.dut.edu.ua зараз недоступний, спробуйте пізніше."

// Holday ..
const Holday = ""

// NoSchedule ..
const NoSchedule = "Заняття відсутні"

// Group ..
const Group = "Група"

// Faculty ..
const Faculty = ""

// Course ..
const Course = ""

// MissingData ..
const MissingData = "Відсутні дані, будь ласка, оберіть свій факультет/курс/групу."

// DefaultMessage ..
func DefaultMessage() string {
	return fmt.Sprintf("Вітаю! я %s. Буду радий вас обслуговувати.\nЯкщо у вас виникнуть питання, будь ласка зверніться до @adambh", Bot.Self.FirstName)
}
