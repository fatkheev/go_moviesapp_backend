package services

import (
    "database/sql/driver"
    "fmt"
    "time"
)

// JSONDate - пользовательский тип для поддержки даты в формате JSON.
type JSONDate time.Time

// IsZero проверяет, является ли дата нулевой.
func (j JSONDate) IsZero() bool {
    return time.Time(j).IsZero()
}

// Value реализует интерфейс driver.Valuer для JSONDate.
func (j JSONDate) Value() (driver.Value, error) {
    return time.Time(j).Format("2006-01-02"), nil
}

// MarshalJSON сериализует дату в JSON.
func (j JSONDate) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))), nil
}

// UnmarshalJSON десериализует дату из JSON.
func (j *JSONDate) UnmarshalJSON(b []byte) error {
    t, err := time.Parse("\"2006-01-02\"", string(b))
    if err != nil {
        return err
    }
    *j = JSONDate(t)
    return nil
}
