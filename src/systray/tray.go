package systray

import (
	"reflect"

	"github.com/energye/systray"
)

type Tray struct {
	title     string
	tooltip   string
	icon      []byte
	menuItems []*MenuItem
	onReady   func()
	onExit    func()
	isRunning bool
}

type MenuItem struct {
	title    string
	tooltip  string
	enabled  bool
	checked  bool
	handler  func()
	systray  *systray.MenuItem
	children []*MenuItem
}

func New(title, tooltip string, icon []byte) *Tray {
	return &Tray{
		title:     title,
		tooltip:   tooltip,
		icon:      icon,
		menuItems: make([]*MenuItem, 0),
	}
}

func (t *Tray) SetIcon(icon []byte) {
	t.icon = icon
	if t.isRunning {
		systray.SetIcon(icon)
	}
}

func (t *Tray) SetTitle(title string) {
	t.title = title
	if t.isRunning {
		systray.SetTitle(title)
	}
}

func (t *Tray) SetTooltip(tooltip string) {
	t.tooltip = tooltip
	if t.isRunning {
		systray.SetTooltip(tooltip)
	}
}

func (t *Tray) AddMenuItem(title, tooltip string, handler func()) *MenuItem {
	item := &MenuItem{
		title:   title,
		tooltip: tooltip,
		enabled: true,
		handler: handler,
	}
	t.menuItems = append(t.menuItems, item)
	return item
}

func (t *Tray) AddSeparator() {
	if t.isRunning {
		systray.AddSeparator()
	}
}

func (m *MenuItem) AddSubMenuItem(title, tooltip string, handler func()) *MenuItem {
	item := &MenuItem{
		title:   title,
		tooltip: tooltip,
		enabled: true,
		handler: handler,
	}
	m.children = append(m.children, item)
	return item
}

func (m *MenuItem) SetEnabled(enabled bool) {
	m.enabled = enabled
	if m.systray != nil {
		if enabled {
			m.systray.Enable()
		} else {
			m.systray.Disable()
		}
	}
}

func (m *MenuItem) SetChecked(checked bool) {
	m.checked = checked
	if m.systray != nil {
		if checked {
			m.systray.Check()
		} else {
			m.systray.Uncheck()
		}
	}
}

func (m *MenuItem) SetTitle(title string) {
	m.title = title
	if m.systray != nil {
		m.systray.SetTitle(title)
	}
}

func (m *MenuItem) SetTooltip(tooltip string) {
	m.tooltip = tooltip
	if m.systray != nil {
		m.systray.SetTooltip(tooltip)
	}
}

func (m *MenuItem) SetHandler(handler func()) {
	m.handler = handler
}

func (t *Tray) OnReady(fn func()) {
	t.onReady = fn
}

func (t *Tray) OnExit(fn func()) {
	t.onExit = fn
}

func (t *Tray) Run() {
	systray.Run(t.setupTray, t.cleanupTray)
}

func (t *Tray) Quit() {
	systray.Quit()
}

func (t *Tray) setupTray() {
	t.isRunning = true
	systray.SetIcon(t.icon)
	systray.SetTitle(t.title)
	systray.SetTooltip(t.tooltip)

	for _, item := range t.menuItems {
		t.createMenuItem(item, nil)
	}

	if t.onReady != nil {
		t.onReady()
	}

	go t.handleEvents()
}

func (t *Tray) createMenuItem(item *MenuItem, parent *systray.MenuItem) {
	var sysItem *systray.MenuItem

	if parent == nil {
		sysItem = systray.AddMenuItem(item.title, item.tooltip)
	} else {
		sysItem = parent.AddSubMenuItem(item.title, item.tooltip)
	}

	item.systray = sysItem

	if !item.enabled {
		sysItem.Disable()
	}
	if item.checked {
		sysItem.Check()
	}

	for _, child := range item.children {
		t.createMenuItem(child, sysItem)
	}
}

func (t *Tray) handleEvents() {
	for _, item := range t.menuItems {
		go t.handleItemEvents(item)
	}
}

func getClickedChan(mi *systray.MenuItem) <-chan struct{} {
	if mi == nil {
		return nil
	}
	v := reflect.ValueOf(mi)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Chan {
			if ch, ok := f.Interface().(<-chan struct{}); ok {
				return ch
			}
		}
	}
	return nil
}

func (t *Tray) handleItemEvents(item *MenuItem) {
	if item.systray == nil {
		return
	}

	if ch := getClickedChan(item.systray); ch != nil {
		for range ch {
			if item.handler != nil {
				item.handler()
			}
		}
	}

	for _, child := range item.children {
		go t.handleItemEvents(child)
	}
}

func (t *Tray) cleanupTray() {
	t.isRunning = false
	if t.onExit != nil {
		t.onExit()
	}
}
