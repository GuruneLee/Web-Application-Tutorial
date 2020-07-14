package main

import (
	"net/http"
	"strings"
)

type router struct {
	//키: http메서드
	//값: URL 패턴별로 실행할 HandlerFunc
	handlers map[string]map[string]http.HandlerFunc
	// map[http메서드] = (map[URL패턴] = handlerFunc )
}

func (r *router) HandleFunc(method, pattern string, h http.HandlerFunc) {
	//http메서드로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]
	if !ok {
		//등록된 맵이 없으면 새 맵을 생성
		m = make(map[string]http.HandlerFunc)
		r.handlers[method] = m
	}
	//http메서드로 등록된 맵에 URL패턴과 핸들러 함수 등록
	m[pattern] = h
}

/*
//정적 URL패턴만 가능한 version.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request){
	if m,ok:=r.handlers[req.Method]; ok{ //http 메서드
		if h,ok:=m[req.URL.Path]; ok{ //URL 패턴
			//요청 URL에 해당하는 핸들러 수행
			h(w, req)
			return
		}
	}
	http.NotFound(w, req)
}
*/

//동적 URL패턴도 가능한 version.
func match(pattern, path string) (bool, map[string]string) {
	//패턴과 패스가 완전동일하면 바로 true반환
	if pattern == path {
		return true, nil
	}

	//패턴과 패스를 "/"단위로 구분
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	//패턴과 패스를 "/"로 구분 후, 문자열 집합의 개수가 다르면 false반환
	if len(patterns) != len(paths) {
		return false, nil
	}

	//패턴에 일치하는 URL매개변수를 담기 위한 params맵 생성
	params := make(map[string]string)

	//"/"로 구분된 패턴,패스의 각 문자열을 하나씩 비교
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
		//부분문자열 일치시 바로 다음 루프 수행
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			//패턴이 ':'로 시작하면 params에 URL params를 담은 후 다음루프 수행
			params[patterns[i][1:]] = paths[i]
		default:
			return false, nil
		}
	}
	//true와 params을 반환
	return true, params
}
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//http메서드에 맞는 모든 handlers를 반복하여 요청
	//요청 URL에 해당하는 handler를 찾는다
	for pattern, handler := range r.handlers[req.Method] {
		if ok, _ := match(pattern, req.URL.Path); ok {
			//요청 URL에 해당하는 handler수행
			handler(w, req)
			return
		}
	}
	//해당 URL에 해당하는 handler를 찾지 못하면 NotFound에러 처리
	http.NotFound(w, req)
	return
}
