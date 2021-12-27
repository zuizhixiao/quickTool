func (m *Template) Create(Db *gorm.DB) error {
    err := Db.Model(&m).Create(&m).Error
    return err
}

func (m *Template) Update(Db *gorm.DB, field ...string) error {
    sql := Db.Model(&m)
    if len(field) > 0 {
        sql = sql.Select(field)
    }
    err := sql.Where("id", m.Id).Updates(m).Error
    return err
}

func (m *Template) GetInfo(Db *gorm.DB) error {
    sql := Db.Model(m).Where("id = ? ", m.Id)
    err := sql.First(&m).Error
    return err
}

func GetTemplateList(Db *gorm.DB, page, num int) ([]Template, error) {
    var list []Template
    sql := Db.Model(Template{})
    if page > 0 && num > 0 {
    sql = sql.Limit(num).Offset((page - 1) * num)
    }
    err := sql.Order("id desc").Find(&list).Error
    return list, err
}

func GetTemplateCount(Db *gorm.DB) (int64, error) {
    var count int64
    sql := Db.Model(Template{})
    err := sql.Count(&count).Error
    return count, err
}
