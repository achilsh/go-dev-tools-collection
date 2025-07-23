package demo

type DemoService interface {
	DoLogic(count float64, id int, name string) (int, error)
}

type ABCServer struct {
	itemer DemoService
}

func (a *ABCServer) ABC(count float64, id int, name string) (int, error) {
	ret, err := a.itemer.DoLogic(count, id, name)
	if err != nil {
		return -1, err
	}
	return ret, nil
}

// 下面是模型被调用时执行额外操作的类型定义：
type User struct {
	ID   int
	Name string
	Age  int
}

type UserService interface {
	UpdateAge(userId int, ageNewVal int) (*User, error)
}

type UserManager struct {
	service UserService
}

func (um *UserManager) Do(userId int, newAge int) (*User, error) {
	return um.service.UpdateAge(userId, newAge)
}
