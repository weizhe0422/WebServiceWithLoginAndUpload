package Utility

func IsValidUser(name string, password string) bool {
	_Name, _Pass, _isValid := "weizhe.chang@gmail.com", "123456", false

	if name == _Name && password == _Pass{
		_isValid = true
	}else{
		_isValid = false
	}
	return _isValid
}
