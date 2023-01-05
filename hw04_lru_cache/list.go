package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head   ListItem
	tail   ListItem
	length int
}

func (l list) Len() int {
	return l.length
}

// Front Первый элемент списка
func (l list) Front() *ListItem {
	return &l.head
}

// Back Последний элемент списка
func (l list) Back() *ListItem {
	return nil
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{v, nil, nil}

	if l.head == (ListItem{}) { // Проверка на пустое значение свойства head. Это ОК ?
		l.tail = item
	} else {
		item.Next = &l.head
		l.head.Prev = &item
	}

	l.head = item
	l.length++
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{v, nil, nil}

	if l.tail == (ListItem{}) { // Проверка на пустое значение свойства head. Это ОК ?
		l.head = item
	} else {
		l.tail.Next = &item
		item.Prev = &l.tail
	}

	l.tail = item // после этого присваивания item.Prev изменяется с последнего значения в списке (10) на себя самого (20, 30) а так же
	l.length++
	return &item
}

func (l list) Remove(i *ListItem) {

}

func (l list) MoveToFront(i *ListItem) {

}

func NewList() List {
	return new(list)
}
