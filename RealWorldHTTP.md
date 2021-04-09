---
title: Real World HTTP
publish-date: 2019-04-19
categories: http
tags:
---

# 개요

* golang으로 서버 파트를 구현하고 curl 로 클라이언트 동작을 테스트하면서 HTTP 를 학습할 수 있는 책. 실습이 많은 점이 좋다.
* HTTP 2.0, QUIC, WebRTC 등의 주제는 그냥 소개만 하는 수준이라 HTTP 1.1 까지의 학습에 적합함.

# Ch01 HTTP/1.0의 신택스: 기본이 되는 네 가지 요소

* 메서드와 경로
* 헤더
* 바디
* 스테이터스 코드

## 실습참고

```bash
go run echo-server.go
```
## curl 사용 샘플

```bash
# GET 메서드로 지정한 data를 urlencode하여 요청
# encode 요청은 브라우저나 curl과 같은 클라이언트에 따라 약간씩 다를 수 있음.
curl --get --data-urlencode "search world" http://localhost:18888

# 서버쪽 로그
# GET /?search%20word HTTP/1.1
# ...

# -v : req/resp 의 헤더를 포함해서 상세한 처리 내용을 보여줌
curl -v http://localhost:18888
```

## 헤더

* 헤더는 '파일명:값' 의 형식. 각 헤더는 한 줄에 하나씩 기술되며 본문과의 사이에는 빈 줄이 하나 있음
* 헤더는 대소문자를 구별하지 않음. 보통은 받는 쪽 라이브러리 나름대로 정규화하여 사용함.
  * 예를 들어 golang의 `net/http`는 `-` 을 구분자로하여 단어를 구분하고 각 단어의 시작은 대문자, 이후는 소문자로 정규화함. `ex. X-TEST -> X-Test`

```bash
curl -H "X-Test: Hello" http://localhost:18888
# RFC상 같은 이름의 헤더를 여러 번 보내는 것도 허용함
# 서버 프로그램에 따라 , 로 구분되는 결합 문자열로 다루거나 배열로 처리하기도 함
curl -H "X-Test: Hello" -H "X-Test: hihi" http://localhost:18888
# golang의 net/http라면 배열로 처리 "X-Test": []string{"Hello", "hihi"}
# python 장고는 {'HTTP_X_TEST': 'Hello,hihi'} 와 같이 정규화 된다고 함

# User-Agent 같이 자주 사용하는 헤더는 --user-agent (-A) 와 같이 alias가 제공됨
# 물론 -H로 직접 지정해도 됨
curl -v -A "Mozilla/5.0"  http://localhost:18888
curl -v -H "User-Agent: Mozilla/5.0"  http://localhost:18888
```

* 헤더는 요청이나 응답에만 사용되는 것, 양쪽 다 사용되는 것이 있음.
* 전체 헤더의 풀 스펙은 https://www.iana.org/assignments/message-headers/message-headers.xhtml 참고

## MIME 타입

* 파일의 종류를 구별하는 문자열, 원래는 전자메일을 위해 만들어짐.
```bash
# 웹 서버가 HTML을 보낼 경우의 응답 헤더 예제
Content-Type: text/html; charset=utf-8
```
* 사진이나 동영상과 같은 미디어는 브라우저, 환경에 따라 이용가능한 포맷이 일부 다름. 때문에 클라이언트와 서버는 다룰 수 있는 포맷에 관해 니고시에이션하고 실제로 반환할 파일 포맷을 변경함.
* content sniffing : IE는 옵션에 따라 MIME 타입이 아닌 내용을 보고 파일 형식을 추측하는 기능이 있음. 서버 설정이 잘못된 경우 이런 추측이 맞으면 장점인 경우도 있겠으나 `text/plain` 인데 스크립트가 들어있다고 브라우저가 멋대로 실행해버리는 경우가 발생하기도 했음. 따라서 보안상의 이유는 아래와 같은 헤더를 전송해 브라우저가 멋대로 MIME 을 추측하지 않도록 하는게 일반적
```
X-Content-Type-Options: nosniff
```
## 전자메일과 HTTP의 비교

* 헤더와 MIME 모두 원래는 전자메일을 위한 기술이었고 이를 HTTP에도 사용하게된 것임.
* 헤더 + 본문 구조는 동일함
* HTTP 요청에는 선두에 '메서드 + 패스' 행이 추가됨
* HTTP 응답에는 선두에 스테이터스 코드가 추가됨
* 그 외에 메일의 경우 긴 헤더가 있을 때 줄바꿈 규칙이 정의되어 있는 등 문법상의 자잘한 차이가 있음
* 그러나 기본은 동일하며 HTTP 통신은 고속으로 전자메일이 왕복하는 것이라고도 볼 수 있음.

## 뉴스그룹

* 메서드와 스테이터스 코드는 뉴스그룹으로부터 도입한것임.
* 즉, 전자메일와 뉴스그룹이 HTTP의 조상

## 메서드

* HTTP는 파일 시스템과 같은 설계 철학으로 만들어짐.
* 가장 흔하게 쓰이는 메서드는 아래 3가지
  * GET : 서버에 헤더와 콘텐츠 요청
  * HEAD : 서버에 헤더만 요청
  * POST : 새로운 문서 투고
* 1.0부터 정의되는 되어 있었지만 브라우저들이 XMLHttpRequest 를 지원하면서부터 사용하게된 메서드
  * PUT : 이미 존재하는 URL의 문서를 갱신
  * DELETE : 지정된 URL의 문서 삭제. 성공하면 삭제된 URL은 무효가 됨.
* 그 외에 1.0이나 1.1 이후로 삭제된 것들. 확실히 파일시스템을 고려한 메서드들이 보임.
  * LINK, UNLINK, CHECKOUT, CHECKIN, SEARCH 등

```bash
# --request= 또는 -X 로 메서드 지정
# 각 메서드별로 단축형도 있음 (ex. HEAD는 --head 혹은 -I)
# 메서드를 생략할 경우 다른 옵션에 따라 GET 이기도 하고 POST 일 때도 있음
## 보통 데이터 전송 옵션을 사용하면 POST가 기본이 됨
curl -X POST http://localhost:18888/greeting
```

## 스테이터스 코드

* 100번대 : 처리가 계속됨을 나타냄.
* 200번대 : 성공했을 때의 응답.
* 300번대 : 서버에서 클라이언트로의 명령. 오류가 아닌 정상처리의 범주. 리다이렉트나 캐시 이용을 지시함.
* 400번대 : 클라이언트가 보낸 요청에 오류가 있음.
* 500번대 : 서버 내부에서 오류가 발생함

## 리다이렉트

* 300번대 응답중 일부는 서버가 브라우저에게 리다이렉트를 지시할 때 사용
* 301 Moved Permanently, 302 Found, 303 See Other, 307 Temporary Redirect, 308 Moved Permanently
* 301, 308 : 요청된 페이지가 영구적으로 이동했을 때 사용. 검색 엔진도 이 응답을 받으면 기존 페이지의 평가를 새로운 페이지로 계승함
* 302, 307 : 일시적인 이동. 모바일 전용 사이트로 이동하거나 관리 페이지 표시 등에 사용.
* 303 : 요청된 페이지에 반환할 컨텐츠가 없거나 원래 반환할 페이지가 따로 있을 때 그쪽으로 이동시키려고 사용. 예를 들면, 로그인 페이지를 사용해 로그인한 후 원래 페이지로 이동하는 경우에 사용.
* 메서드 변경 : 첫 번째 요청이 POST이고, 두 번째 이후에 GET, HEAD를 사용할 경우 사용자에게 확이할 필요 없이 실시할 수 있는지여부
  * 301, 302는 허용, 303, 307, 308은 허가 필요
* 영구적/일시적
  * 301, 303, 308은 영구적, 302, 307은 일시적
* 캐시 
  * 301, 308은 함. 302, 307은 지시에 따름 (Cache-Control, Expires 헤더 등), 303은 캐싱 안함.

* 클라이언트는 Location 헤더를 보고 재전송함. 재전송할 때는 헤더 등도 다시 보냄.
* curl은 -L 옵션을 부여하면, 응답이 300번대고 Location 헤더가 있으면 Location 헤더에서 지정한 URL에 재전송을 수행함.

```bash
# 기본적으로 최대 50번까지 리다이렉트, --max-redirs 옵션으로 지정가능
curl -L http://localhost:18888
```

* 리다이렉트 횟수는 스펙상 정해진 제한은 없고 클라이언트가 리다이렉트 무한을 탐지해야함.
  * curl은 기본 50회 제한, golang에서는 10회 제한 등 구현에 따라 다름.
* 구글 권장은 5회 이하, 가능하면 3회 이하임.

## URL

* URL (Uniform Resource Locator) : 장소로 문서 등의 리소스를 특정하는 수단. 즉, 주소.
* URN (Uniform Resource Name) : `urn:ietf:rfc:1738` 과 같이 이름 그자체. 이름밖에 없으므로 실제 위치를 알려면 따로 정보가 필요.
* URI : URL, URN을 포함하는 개념. 웹에서 URN이 사용될 일이 없으므로 URL과 URI는 거의 같음.
* RFC 3305에서는 URL은 관용 표현, URI는 공식 표기로 정의했으나 실제로는 URL이 더 일반적으로 널리 쓰임. 언어에 따라서는 Ruby, C#은 URI 라 하는 반면 golang, python은 URL 을 사용함. 근데 JAVA는 URI, URL 클래스가 다 있음.

#### URL의 구조

* `스키마://호스트명/경로?쿼리` : 일반적인 형식
  * https://www.oreilly.co.kr/index.html?q=123
* `스키마://사용자:패스워드@호스트명:포트/경로#프래그먼트?쿼리` : 모든 요소를 포함한 경우
* 스키마의 해석은 브라우저와 같은 클라이언트의 책임
* 사용자, 패스워드 : Basic 인증 방식에 사용
* 프래그먼트 : HTML 페이지 내 링크의 앵커를 지정하는 데 사용
* 인코딩
  * URL은 기본적으로 ASCII 문자열로, 영문자, 숫자, 몇 개의 기호만 표시할 수 있음
  * RFC 2718 부터는 UTF-8로 URL을 인코딩하므로 다국어 문자도 다룰 수 있음.
* 스펙상 URL의 길이 제한은 없음. 근데 IE는 2083자까지 다를 수 있어 대체로 2000자 정도를 제한으로 봄. (내 의견 : IE 제외한 모던 브라우저들은 다를 지도. 그러나 애초에 URL이 2000자를 넘을정도면 URL 설계를 잘못했다고 본다.)
* HTTP2 에서는 URL이 지나치게 길 때 반환되는 스테이터스 코드 414 URI Too Long 이 추가됨.

#### URL과 국제화

* 퓨니코드를 사용하면 도메인 이름에 한글, 한자같은 문자를 다룰 수 있다.
* 실제로 UTF-8 같은 인코딩을 쓰는건 아니고 퓨니코드에서 지정한 규칙으로 변환하는 방식을 사용한다.
* 예를 들어 `한글도메인.kr` 을 퓨니코드로하면 `xn--bj0bj3i97fq8o5lq.kr` 이다.

## 바디 (Body)

#### 응답의 바디

* 읽어 올 바이트 수는 Content-Length 헤더로 지정한다.
  * Content-Encoding을 통해 압축된 경우 압축 후의 크기를 의미한다.
* HEAD 메서드 요청일 때도 Content-Length 헤더를 반환해야한다. 캐시용 ETag 등도 마찬가지.

#### 요청의 바디

```bash
# -d, --data, --data-ascii : 텍스트 데이터 (인코딩 없이 그대로 보냄)
# --data-urlencode : 텍스트 데이터 (curl 커맨드가 보내기 전 인코딩 수행)
# --data-binary : 바이너리 데이터
# -T 파일명 혹은 -d @파일명 : 보내고 싶은 데이터를 파일에서 읽음
# -d, --data-urlencode 는 기본적으로 application/x-www-form-urlencoded 으로 보냄

curl -d "name=bob" http://localhost:18888
curl -d "{\"hello\": \"world\"}" -H "Content-Type: application/json" http://localhost:18888
# test.json 
curl -d @test.json -H "Content-Type: application/json"  http://localhost:18888
```

* GET, HEAD, DELETE, OPTIONS, CONNECT 는 페이로드 바디를 가질 수는 있지만, 구현에 따라서는 서버가 이를 받아들이지 않을 수 있음 (RFC 7231)
* TRACE 는 "페이로드 바디를 포함해선 안 된다" 라고 강조되어있음

> 내의견 : 제대로된 요청이라면 GET 과 같은 메서드에 바디를 담아 보내서는 안 된다.

## 정리

* 메서드와 경로, 헤더, 바디, 스테이터스 코드는 HTTP의 기초
* 이는 HTTP/2 에서도 바뀌지 않았음.