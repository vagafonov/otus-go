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
	head   *ListItem
	tail   *ListItem
	length int
}

func (l list) Len() int {
	return l.length
}

// Front Первый элемент списка.
func (l list) Front() *ListItem {
	return l.head
}

// Back Последний элемент списка.
func (l list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{v, nil, nil}

	if l.head == nil {
		l.head = &item
		l.tail = &item
	} else {
		item.Next = l.head
		l.head.Prev = &item
	}

	l.head = &item
	l.length++
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	/*
		TODO ?
		Является ли объявление структуры с "&" в начале &ListItem{v, nil, nil}
		и без дальнейшего использования взятия указателя "&" от переменной
		тем же самым что реализовано ниже. Если да то какой вариант предпочтительнее ?
	*/

	item := ListItem{v, nil, nil}

	if l.tail == nil { // Проверка на пустое значение свойства head. Это ОК ?
		l.head = &item
	} else {
		item.Prev = l.tail
		l.tail.Next = &item
	}

	l.tail = &item
	l.length++
	return &item
}

func (l *list) Remove(i *ListItem) {
	switch {
	// Удаление не первого и не последнего элемента
	case i.Prev != nil && i.Next != nil:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	// Удаление первого элемента
	case i.Prev == nil && i.Next != nil:
		i.Next.Prev = nil
		// TODO ? Мало удалить сам элемент (i = nil) кроме этого нужно удалить все сылки на него l.head, i.Next.Prev ?
		l.head = i.Next
	// Удаление последнего элемента
	case i.Next == nil && i.Prev != nil:
		i.Prev.Next = nil
		// TODO ? Мало удалить сам элемент (i = nil) кроме этого нужно удалить все сылки на него l.head, i.Next.Prev ?
		l.tail = i.Prev
	default:
		i = nil
	}

	i = nil
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	// Перемещение последнего в конец
	if i.Next == nil {
		// Последний элемент в списке не один
		if i.Prev != nil {
			i.Prev.Next = nil
			i.Prev = nil
			i.Next = l.head
			l.head = i
		}
	}
	// Перемещение первого в начало бессмысленно
}

func NewList() List {
	return new(list)
}
