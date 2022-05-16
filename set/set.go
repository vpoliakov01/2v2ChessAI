package set

type Set struct {
	set map[interface{}]struct{}
}

func New(elems ...interface{}) *Set {
	m := map[interface{}]struct{}{}
	s := &Set{set: m}

	for i := range elems {
		s.Add(elems[i])
	}

	return s
}

func (s *Set) Add(elem interface{}) {
	s.set[elem] = struct{}{}
}

func (s *Set) Delete(elem interface{}) {
	delete(s.set, elem)
}

func (s *Set) Has(elem interface{}) bool {
	_, has := s.set[elem]
	return has
}

func (s *Set) Elements() []interface{} {
	res := []interface{}{}
	for elem := range s.set {
		res = append(res, elem)
	}
	return res
}

func (s *Set) Clear() {
	s.set = map[interface{}]struct{}{}
}

func (s *Set) Size() int {
	return len(s.set)
}

func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

func (s Set) Copy() *Set {
	return New(s.Elements()...)
}
