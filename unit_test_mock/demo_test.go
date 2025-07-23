package demo

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 定义 mock 依赖的对象，该对象实现 依赖的方法。
type mockDemoService struct {
	mock.Mock
}

// 定义mock 依赖对象的 内部依赖方法
func (m *mockDemoService) DoLogic(count float64, id int, name string) (int, error) {
	// 记录实际调用，并返回预设的返回值（由 On().Return() 定义）
	args := m.Called(count, id, name)

	// 从返回值中解析结果
	return args.Int(0), args.Error(1)
}

// 测试 mock依赖的正确处理逻辑
func TestMockDemoService(t *testing.T) {

	// 初始化模拟对象
	mserver := new(mockDemoService)

	// mock 对应方法， 入参，出参； Return是设置返回值； 预设-设置预期行为：当调用 DoLogic(12.101, 11, "abc") 时，返回预设 100 和 nil 错误
	mserver.On("DoLogic", 12.101, 11, "abc").Return(100, nil)

	// 初始化待测试服务（注入模拟对象)
	abcItem := &ABCServer{
		itemer: mserver,
	}

	//调用上层逻辑和方法； 执行测试方法
	ret, err := abcItem.ABC(12.101, 11, "abc")
	assert.NoError(t, err)
	assert.Equal(t, ret, 100)

	// 5. 验证 mock 对象的预期调用是否发生
	mserver.AssertExpectations(t)
}

// 测试 mock 依赖的异常处理逻辑
func TestMockErrorDemoService(t *testing.T) {
	// 定义mock 对象
	mserver := new(mockDemoService)
	// mock 对应方法， 入参，出参
	mockError := fmt.Errorf("is new error")
	mserver.On("DoLogic", 12.101, 11, "abc").Return(-1, mockError)

	abcItem := &ABCServer{
		itemer: mserver,
	}

	//调用上层逻辑和方法
	ret, err := abcItem.ABC(12.101, 11, "abc")
	assert.ErrorIs(t, err, mockError)
	assert.Equal(t, ret, -1)

	// 5. 验证 mock 对象的预期调用是否发生
	mserver.AssertExpectations(t)
}

// 设置 mock 方法 的 自定义逻辑，允许在模拟被调用时执行额外操作，比如：修改参数，抛出异常，打印日志等。
type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) UpdateAge(userId int, ageNewVal int) (*User, error) {
	args := m.Called(userId, ageNewVal)
	return args.Get(0).(*User), args.Error(1)
}

func TestMockRunDo(t *testing.T) {
	m := new(mockUserService)
	c := m.On("UpdateAge", mock.AnythingOfType("int"), mock.AnythingOfType("int"))

	retErr := errors.New("invalid age: cannot be negative")
	fn := func(args mock.Arguments) {
		userID := args.Get(0).(int)
		newAge := args.Get(1).(int)

		// 自定义逻辑1：检查年龄合法性
		if newAge < 0 {
			// 设置返回值为 nil 和错误
			c.Return((*User)(nil), retErr)
			return
		}

		// 自定义逻辑2：构造返回对象

		c.Return(&User{
			ID:   userID,
			Name: "test_user",
			Age:  newAge,
		}, nil)
	}
	// 定义 mocK 方法的预期和自定义逻辑；  使用 Run 方法设置自定义逻辑；返回 Return() 预设的返回值（若 Run() 中动态修改了返回值，以修改后为准）
	c.Run(fn).Return((*User)(nil), error(nil))

	manager := &UserManager{
		service: m,
	}

	t.Run("valid age", func(t *testing.T) {
		u, err := manager.Do(100, 100)
		assert.NoError(t, err)
		assert.Equal(t, 100, u.Age)
		assert.Equal(t, 100, u.Age)
	})

	t.Run("invalid age", func(t *testing.T) {
		u, err := manager.Do(0, -90)
		assert.Error(t, err)
		assert.Nil(t, u)
		assert.ErrorIs(t, err, retErr)

	})
}

// 1. 定义接口与数据结构
type Person struct {
	ID        int
	Name      string
	UpdatedAt time.Time
}
type PersonRepository interface {
	UpdatePerson(user *Person) error
}

// 2. 实现模拟对象
type MockPersonRepository struct {
	mock.Mock
}

func (m *MockPersonRepository) UpdatePerson(user *Person) error {
	args := m.Called(user)
	return args.Error(0)
}

// 3. 测试用例：验证 UpdateUser 被调用时，UpdatedAt 被正确设置
func TestUpdateUser_SideEffect(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	person := &Person{ID: 1, Name: "Bob"}

	// 预设预期：当调用 UpdateUser(user) 时，执行 Run() 中的逻辑，再返回 nil 错误
	mockRepo.On("UpdatePerson", person).
		Run(func(args mock.Arguments) {

			// 从 args 中获取实际传入的 user 对象（注意类型断言）
			actualUser := args[0].(*Person)

			// 模拟副作用：设置 UpdatedAt 为当前时间
			actualUser.UpdatedAt = time.Now().Truncate(time.Second)
		}).
		Return(nil) // 预设返回值

	// 调用模拟方法
	err := mockRepo.UpdatePerson(person)

	// 验证结果
	assert.NoError(t, err)
	assert.NotEqual(t, time.Time{}, person.UpdatedAt) // 确认副作用生效
	mockRepo.AssertExpectations(t)
}

// 1. 定义接口
type Calculator interface {
	Add(a, b int) int
}

// 2. 模拟对象
type MockCalculator struct {
	mock.Mock
}

func (m *MockCalculator) Add(a, b int) int {
	args := m.Called(a, b)
	return args.Int(0) // 返回 Call 中预设的返回值（动态生成的 result）
}

// 3. 测试用例：动态计算返回值
func TestCalculator_DynamicReturn(t *testing.T) {
	mockCalc := new(MockCalculator)

	// 定义一个变量存储动态计算的结果（通过闭包传递给 Return）
	var dynamicResult int

	// 预设：调用 Add(a, b) 时，通过 Run() 计算 a+b 保存 *Call 对象：mockCalc.On(...) 返回的 *Call 对象被赋值给 call，用于后续动态修改返回值。
	c := mockCalc.On("Add", mock.Anything, mock.Anything)
	//通过 call.Return(result) 直接将结果设置为该调用的返回值；这会覆盖 On() 方法默认的返回值（无需再显式调用 Return()）
	c.Run(func(args mock.Arguments) {
		a := args[0].(int)
		b := args[1].(int)
		dynamicResult = a + b   // 计算结果存入外部变量
		c.Return(dynamicResult) // 关键：用计算出的 result 覆盖默认返回值
	}).Return(0)

	// 调用模拟方法（触发 Run() 中的计算）
	sum := mockCalc.Add(3, 5)

	// 验证结果：此时 dynamicResult 已被 Run() 赋值为 8，模拟方法返回 8
	assert.Equal(t, 8, sum)
	mockCalc.AssertExpectations(t)
}
