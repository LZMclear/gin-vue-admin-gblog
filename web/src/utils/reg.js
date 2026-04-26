export const checkEmail = (rule, value, callback) => {
  if (!value) {
    callback()
    return
  }
  const valid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
  valid ? callback() : callback(new Error('邮箱格式不正确'))
}
