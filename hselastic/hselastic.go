package hselastic

type QueryBool struct {
	Query				`json:"query,omitempty"`
}

type Query struct {
	Bool				`json:"bool,omitempty"`
}

func HsESQuery()(*QueryBool){
	q := new(QueryBool)

	q.Bool.Should = make(should,0,0)
	q.Bool.Must = make(must,0,0)
	q.Bool.Mustnot = make(mustnot,0,0)
	q.Bool.Term = make(term,0,0)
	q.Bool.Filter = nil

	return q
}

func (s *QueryBool)Should()(*should){
	return &s.Bool.Should
}

func (s *QueryBool)Must()(*must){
	return &s.Bool.Must
}

func (s *QueryBool)Mustnot()(*mustnot){
	return &s.Bool.Mustnot
}

func (s *QueryBool)Term()(*term){
	return &s.Bool.Term
}

func (s *QueryBool)Filter()(*setRange){
	if s.Bool.Filter == nil{
		ran := new(setRange)
		ran.Setrange = make(map[string](map[string]interface{}))
		s.Bool.Filter = ran
	}
	return s.Bool.Filter
}

type must 		[]*match
type mustnot 	[]*match
type should 	[]*match
type term 		[]*match
type filter 	*setRange

/*ES querybool should 处理*/
type Bool struct {
	Must 	must 			`json:"must,omitempty"`
	Mustnot mustnot  		`json:"must_not,omitempty"`
	Should	should       	`json:"should,omitempty"`
	Term 	term   			`json:"term,omitempty"`
	Filter 	filter			`json:"filter,omitempty"`
}

/*ES querybool should 处理*/
func (s *should)Match()(* match){
	m := new(match)

	*s = append(*s, m)
	return m
}

/*ES querybool must 处理*/
func (mu *must)Match()(* match){
	m := new(match)

	*mu = append(*mu, m)
	return m
}

/*ES querybool must_not 处理*/
func (mu *mustnot)Match()(* match){
	m := new(match)

	*mu = append(*mu, m)
	return m
}

/*ES querybool term 处理*/
func (t *term)Match()(* match){
	m := new(match)

	*t = append(*t, m)
	return m
}

/*ES querybool Filter 处理*/
/*ES 设置过滤器范围查询条件*/
type setRange struct {
	Setrange map[string](map[string]interface{})		`json:"range,omitempty"`
}

func (m *setRange)setRange(k,c string, v interface{}){
	if v == nil{
		return
	}

	mp := make(map[string]interface{})
	mp[c] = v
	m.Setrange[k] = mp
}

/*ES 匹配对象查询*/
type match struct {
	Match interface{}		`json:"match,omitempty"`
}

func (m *match)Search(k string, v interface{}){
	if v == nil{
		return
	}

	var mp = make(map[string]interface{},1)
	mp[k] = v
	m.Match = mp
}