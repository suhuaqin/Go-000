# 学习笔记

## 关于 std errors 库和 github/pkg/errors 库

### 携带堆栈信息
- github/pkg/errors 支持，使用errors.WithStack()与errors.Wrap()， 使用 %+v 输出
- std errors 需要自己实现
### 携带消息
- github/pkg/errors  使用errors.WithMessage() 携带消息
- std errors 使用 fmt.Errorf("%w ...", err, ...) 携带消息, 与 pkg/errors 的 errors.WithMessagef() 等效
### 取出被包装的错误
- 两个库都可以使用 errors.Unwrap() 获取，是兼容的
### 与特定错误的比较
- 使用Errors.Is()比较，两个errors库都可以这样用，是兼容的