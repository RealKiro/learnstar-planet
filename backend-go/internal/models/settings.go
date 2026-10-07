// 设置读写（settings JSON 文本列）：用户级 / 学校级 / 班级级。
// 语义忠实移植自 Laravel App\Models\User::getSetting / setSetting：
// 键不存在时返回调用方给出的默认值，写入时按 key 合并（不清空其他键）。
package models

import "encoding/json"

// parseSettingsMap 解析 settings JSON 文本；为空或解析失败时返回空表。
func parseSettingsMap(raw string) map[string]any {
	m := map[string]any{}
	if raw == "" {
		return m
	}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return map[string]any{}
	}
	return m
}

// readSettingString 读取键值表中的字符串设置，缺省/非字符串/空串返回 def。
func readSettingString(m map[string]any, key, def string) string {
	v, ok := m[key]
	if !ok {
		return def
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return def
	}
	return s
}

// readSettingUint 读取键值表中的正整数设置，未设置/非法/为 0 时返回 0。
func readSettingUint(m map[string]any, key string) uint {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		if n <= 0 {
			return 0
		}
		return uint(n)
	case int:
		if n <= 0 {
			return 0
		}
		return uint(n)
	}
	return 0
}

// mergeSetting 返回「合并该键后」的 settings JSON 文本（不落库，由调用方保存）。
func mergeSetting(m map[string]any, key string, value any) (string, error) {
	m[key] = value
	raw, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// SettingsMap 解析用户设置为键值表；为空或解析失败时返回空表。
func (u *User) SettingsMap() map[string]any {
	if u == nil {
		return map[string]any{}
	}
	return parseSettingsMap(u.Settings)
}

// SettingString 读取字符串设置，缺省返回 def。
func (u *User) SettingString(key, def string) string {
	return readSettingString(u.SettingsMap(), key, def)
}

// SettingUint 读取正整数设置，未设置/非法/为 0 时返回 0。
func (u *User) SettingUint(key string) uint {
	return readSettingUint(u.SettingsMap(), key)
}

// WithSetting 返回「合并该键后」的 settings JSON 文本（不落库，由调用方保存）。
func (u *User) WithSetting(key string, value any) (string, error) {
	return mergeSetting(u.SettingsMap(), key, value)
}

// SettingsMap 解析学校设置为键值表（Laravel schools.settings）。
func (s *School) SettingsMap() map[string]any {
	if s == nil {
		return map[string]any{}
	}
	return parseSettingsMap(s.Settings)
}

// SettingString 读取学校字符串设置，缺省返回 def。
func (s *School) SettingString(key, def string) string {
	return readSettingString(s.SettingsMap(), key, def)
}

// WithSetting 返回「合并该键后」的学校 settings JSON 文本（不落库，由调用方保存）。
func (s *School) WithSetting(key string, value any) (string, error) {
	return mergeSetting(s.SettingsMap(), key, value)
}

// SettingsMap 解析班级设置为键值表（Laravel class_rooms.settings）。
func (c *ClassRoom) SettingsMap() map[string]any {
	if c == nil {
		return map[string]any{}
	}
	return parseSettingsMap(c.Settings)
}

// SettingString 读取班级字符串设置，缺省返回 def。
func (c *ClassRoom) SettingString(key, def string) string {
	return readSettingString(c.SettingsMap(), key, def)
}

// WithSetting 返回「合并该键后」的班级 settings JSON 文本（不落库，由调用方保存）。
func (c *ClassRoom) WithSetting(key string, value any) (string, error) {
	return mergeSetting(c.SettingsMap(), key, value)
}
