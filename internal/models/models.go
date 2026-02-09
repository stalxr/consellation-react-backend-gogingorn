package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User - модель пользователя системы (донор или администратор)
// Хранит информацию о пользователе: email, пароль (хэшированный), имя, роль и общую сумму пожертвований
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"column:password_hash;not null" json:"-"` // Пароль не возвращается в JSON
	FullName     string         `gorm:"column:full_name;not null" json:"full_name"`
	AvatarURL    *string        `gorm:"column:avatar_url" json:"avatar_url,omitempty"`
	Role         string         `gorm:"type:varchar(20);default:'user';not null" json:"role"` // Роль: "user" (донор) или "admin" (администратор)
	TotalDonated float64        `gorm:"column:total_donated;default:0" json:"total_donated"` // Общая сумма всех пожертвований пользователя
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Мягкое удаление
}

// TableName возвращает имя таблицы в базе данных для модели User
func (User) TableName() string {
	return "users"
}

// BeforeCreate генерирует UUID для пользователя перед сохранением в БД.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Dream - модель мечты (благотворительного сбора)
// Представляет карточку мечты с информацией о цели сбора, собранной сумме и статусе
type Dream struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title           string         `gorm:"not null" json:"title"`                                    // Название мечты
	Slug            string         `gorm:"uniqueIndex" json:"slug"`                                  // Уникальный идентификатор для URL (например: "vanya-more")
	ShortDescription *string       `gorm:"column:short_description" json:"short_description,omitempty"` // Краткое описание
	FullDescription *string       `gorm:"column:full_description;type:text" json:"full_description,omitempty"` // Полное описание (HTML контент)
	TargetAmount    float64        `gorm:"column:target_amount;not null" json:"target_amount"`      // Целевая сумма сбора
	CollectedAmount float64        `gorm:"column:collected_amount;default:0" json:"collected_amount"` // Собранная сумма
	Status          string         `gorm:"type:varchar(20);default:'active';not null" json:"status"` // Статус: "active", "completed", "frozen"
	CoverImage      *string        `gorm:"column:cover_image" json:"cover_image,omitempty"`        // URL обложки
	GalleryImages   []string       `gorm:"serializer:json" json:"gallery_images,omitempty"`             // Массив URL изображений галереи
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"-"`
	ClosedAt        *time.Time     `gorm:"column:closed_at" json:"closed_at,omitempty"`           // Дата закрытия сбора
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`                                         // Мягкое удаление
}

// TableName возвращает имя таблицы в базе данных для модели Dream
func (Dream) TableName() string {
	return "dreams"
}

// BeforeCreate генерирует UUID для мечты перед сохранением в БД.
func (d *Dream) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// Donation - модель пожертвования
// Хранит информацию о пожертвовании: сумму, связь с мечтой и пользователем, статус оплаты
type Donation struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	DreamID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"dream_id"`              // ID мечты, на которую пожертвовано
	UserID      *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`              // ID пользователя (nullable - может быть анонимным)
	Amount      float64        `gorm:"not null" json:"amount"`                                 // Сумма пожертвования
	Email       *string        `gorm:"column:email" json:"email,omitempty"`                  // Email для анонимных пожертвований
	IsAnonymous bool           `gorm:"column:is_anonymous;default:false" json:"is_anonymous"` // Флаг анонимности
	Comment     *string        `gorm:"type:text" json:"comment,omitempty"`                     // Комментарий донора
	Status      string         `gorm:"type:varchar(20);default:'pending';not null" json:"status"` // Статус: "pending", "completed", "failed"
	PaymentURL  *string        `gorm:"column:payment_url" json:"payment_url,omitempty"`       // Ссылка на страницу оплаты (эквайринг)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"` // Мягкое удаление

	// Связи с другими моделями (не возвращаются в JSON)
	Dream Dream `gorm:"foreignKey:DreamID" json:"-"`
	User  *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName возвращает имя таблицы в базе данных для модели Donation
func (Donation) TableName() string {
	return "donations"
}

// BeforeCreate генерирует UUID для пожертвования перед сохранением в БД.
func (d *Donation) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// News - модель новости фонда
// Хранит информацию о новостях: заголовок, текст, изображение и дату публикации
type News struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`                        // Заголовок новости
	PreviewText *string        `gorm:"column:preview_text;type:text" json:"preview_text,omitempty"` // Краткий текст для превью
	Content     *string        `gorm:"type:text" json:"content,omitempty"`            // Полный текст новости
	ImageURL    *string        `gorm:"column:image_url" json:"image_url,omitempty"`   // URL изображения новости
	PublishedAt *time.Time     `gorm:"column:published_at" json:"published_at,omitempty"` // Дата публикации
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"` // Мягкое удаление
}

// TableName возвращает имя таблицы в базе данных для модели News
func (News) TableName() string {
	return "news"
}

// BeforeCreate генерирует UUID для новости перед сохранением в БД.
func (n *News) BeforeCreate(tx *gorm.DB) (err error) {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// Report - модель финансового отчета фонда
// Хранит информацию о финансовых отчетах: год, месяц, название и ссылку на PDF файл
type Report struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Year      int            `gorm:"not null" json:"year"`                    // Год отчета
	Month     int            `gorm:"not null" json:"month"`                   // Месяц отчета (1-12)
	FileURL   string         `gorm:"column:file_url;not null" json:"file_url"` // Ссылка на PDF файл отчета
	Title     string         `gorm:"not null" json:"title"`                  // Название отчета
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"` // Мягкое удаление
}

// TableName возвращает имя таблицы в базе данных для модели Report
func (Report) TableName() string {
	return "reports"
}

// BeforeCreate генерирует UUID для отчёта перед сохранением в БД.
func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

