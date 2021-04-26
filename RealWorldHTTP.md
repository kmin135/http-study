---
title: Real World HTTP
publish-date: 2019-04-19
categories: http
tags:
---

# 개요

* golang으로 서버 파트를 구현하고 curl 로 클라이언트 동작을 테스트하면서 HTTP 를 학습할 수 있는 책. 실습이 많은 점이 좋다.
* HTTP 2.0, QUIC, WebRTC 등의 주제는 그냥 소개만 하는 수준이라 HTTP 1.1 까지의 학습에 적합함.
* https://developer.mozilla.org/en-US/docs/Web/HTTP 도 정리가 잘 되어있다. 브라우저별 지원여부나 관련 RFC도 잘 정리 되어있어서 추후 레퍼런스로 참고하기에 매우 적합해보인다.

# Ch01 HTTP/1.0의 신택스: 기본이 되는 네 가지 요소

* 메서드와 경로
* 헤더
* 바디
* 스테이터스 코드

## 실습참고

```bash
go run echo-server.go
```

## http 기본 구조

```bash
# 헤더와 바디는 빈 개행으로 구분한다.
# 요청
[메서드] [경로] [프로토콜]
[헤더]

[바디]

# 요청 ex
POST / HTTP/1.1
Host: localhost:18888
User-Agent: curl/7.68.0

title=The%20%26%20Art&author=Bob

# 응답
[프로토콜] [스테이터스 코드 + 메시지]
[헤더]

[바디]

# 응답 ex
HTTP/1.1 200 OK
Content-Length: 32
Content-Type: text/html; charset=utf-8

<html><body>hello</body></html>
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

# Ch02 HTTP/1.0의 시맨틱스: 브라우저 기본 기능의 이면

## 기본 form 전송 (x-www-form-urlencoded)

* form의 기본 전송 MIME 타입
* RFC에 x-www-form-urlencoded 타입인 경우의 파일전송동작은 정의되지 않았음.
  * 브라우저에서 해보면 파일명만 전달되고 파일은 전송되지 않음.

```html
<!--
- 브라우저는 RFC 1866에서 책정한 변환 포맷으로 변한한다.
- 알파벳, 숫자, 별(*), 하이픈(-), 마침표(.), 언더스코어(_) 의 여섯 종류 문자 외에는 변환 필요
- 예를 들어 공백은 + 로 바뀜
- 값으로 들어가는 =은 %3D, &는 %26으로 바뀌고 실제 구분자는 =, & 로 전달되므로 읽는 쪽에서 이를 구분할 수 있음

title=The & Art, author=Bob으로 요청한 경우의 바디
title=The+%26+Art&author=Bob
-->
<form action="http://localhost:18888" method="POST">
  <input name="title"/>
  <input name="author"/>
  <input type="submit"/>
</form>

<!-- 
 - method가 GET일 경우 바디가 아니라 쿼리로서 URL에 부여함 (RFC 1866 정의)
 GET /?title=The+%26+Art&author=Bob HTTP/1.1
-->
<form action="http://localhost:18888" method="GET">
<!-- 이하 동일 -->

<!--
쓸 일은 없겠으나 form은 text/plain 타입도 지원한다.
변환을 하지않으며, 개행으로 구분해 값을 전송한다.

title=The & Art
author=Bob
-->
<form action="http://localhost:18888" method="POST" enctype="text/plain">
<!-- 이하 동일 -->
```
```bash
# -d 는 인코딩을 하지 않고 그대로 보내므로 아래 예의 "The & Art" 같은 인코딩이 필요한건 보내면 안 된다.
curl -d title="The & Art" -d author="Bob" http://localhost:18888

# --data-urlencode 옵션은 브라우저와 유사하게 인코딩해서 보냄. 단, 브라우저와 달리 RFC 3986 에서 정의한 변환 방식을 사용함
# 예를 들어 공백이 + 가 아닌 %20 으로 변환됨

# title=The%20%26%20Art&author=Bob
curl --data-urlencode title="The & Art" --data-urlencode author="Bob" http://localhost:18888
```

* 웹 브라우저는 form 인코딩에 RFC 1866, curl은 데이터 인코딩에 RFC 3986 을 사용하지만 동일 알고리즘으로 복원할 수 있음

## form을 이용한 파일 전송 (multipart/form-data)

* RFC 1867에서 정의
* HTTP 바디는 일반적으로 한 파일 전체를 의미하고 단순히 Content-Length 만큼 읽으면 됨.
* 반면에 멀티파트는 이름 그대로 한 번의 요청으로 복수의 파일을 전송할 수 있음. 받는 쪽에서 복수의 파일을 구분하기 위한 방법이 필요함. 이를 위해 Content-Type 헤더에 boundary 라는 경계문자열을 부여함.

```html
<!-- 
- Content-Type에 부여되는 boundary 는 각 클라이언트가 랜덤으로 생성함
- boundary 값으로 body의 데이터를 분리하여 해석할 수 있음
- body 끝에는 [boundary값]-- 으로 끝남
- boundary로 구분되는 각 Part는 각각 헤더+빈줄+콘텐츠로 구성됨
- 각 파트에 Content-Disposition 헤더가 있는데 이를 통해 각 파트를 정의함
  - 내 의견) Disposition 은 기질, 성향, 배치 라는 뜻. 책 번역은 기질,성질이라고 해놨는데 배치(arrangement) 가 더 적절한 해석이라고 생각한다.
- 파일을 첨부해보면 Content-Type 헤더가 부여됨을 알 수 있음

[헤더들]
Content-Type: multipart/form-data; boundary=---------------------------340904805324056899591929825476
[나머지헤더들]

[이하 바디]

-----------------------------340904805324056899591929825476
Content-Disposition: form-data; name="title"

The & Art
-----------------------------340904805324056899591929825476
Content-Disposition: form-data; name="author"

Bob
-----------------------------340904805324056899591929825476
Content-Disposition: form-data; name="attachment"; filename="sample.txt"
Content-Type: text/plain

hello file
-----------------------------340904805324056899591929825476--  
-->

<form action="http://localhost:18888" method="POST" enctype="multipart/form-data">
  <input name="title"/>
  <input name="author"/>
  <input name="attachment" type="file">
  <input type="submit"/>
</form>
```

```bash
# curl은 -d 대신 -F 를 사용하면 multipart/form-data 로 전송함
# -d, -F는 함께 사용할 수 없음
# type, filename은 생략가능함. 이 경우 type은 자동 설정되고 filename은 로컬 파일명과 동일
: '
Content-Disposition: form-data; name="attachment"; filename="changed.txt"
Content-Type: application/json

{"hello": "world"}
'
curl -F "title=The & Art" -F author=Bob -F "attachment=@test.json;type=application/json;filename=changed.txt" http://localhost:18888

# =<[파일] 형식으로 파일을 첨부하는게 아니라 파일 내용을  보내는것도 가능
: '
Content-Disposition: form-data; name="attachment"
Content-Type: application/json

{"hello": "world"}
'
curl -F "title=The & Art" -F author=Bob -F "attachment=<test.json;type=application/json" http://localhost:18888

```

## form 을 이용한 리다이렉트

* 보통은 300번대 스테이터스 코드를 이용해서 리다이렉트를 수행함.
* 이 방법의 제한사항은
  * URL은 환경에 따라 길이 제한이 있을 수 있으므로 GET의 쿼리로 보낼 수 있는 데이터 양에는 한계가 있음
  * 데이터가 URL에 포함되므로 전송 내용이 액세스 로그 등에 남을 수 있음
* 이런 경우 form 을 이용한 리다이렉트를 사용할 수 있음

```html
<!-- body의 onload 이벤트를 활용한 간단한 구조 -->
<body onload="document.forms[0].submit()">
  <form action="http://localhost:18888/redirect" method="POST">
    <input type="hidden" name="data" value="i want to send this data" />
  </form>
</body>
```

* 전송가능한 데이터양에 제한이 없으므로 리다이렉트 전송할 데이터가 많을 때 유용함
* 순간적으로 빈 페이지가 표시될 수 있다는 게 단점.
> 내 의견 : 이라고 했지만 굳이 onload 이벤트를 안 해도 중간 페이지를 별도로 디자인하고 화면이 다 뜨고 적절한 UI를 보여준 후 리다이렉트해도 되므로 단점이라 볼 건 없는듯.
* 또 300번대 스테이터스 코드와 달리 클라이언트 환경에서 자바스크립트가 비활성화 되어있으면 자동 전환이 동작하지 않는것도 차이점.
* SOAP 형식의 조금 큰 XML 데이터를 암호화하여 리다이렉트할 필요가 있는 SAML, OpenID Connect 등에서 활용됨.

## 콘텐트 니고시에이션

* 통신 방법을 최적화하고자 하나의 요청 안에서 서버와 클라이언트가 서로 최고의 설정을 공유하는 시스템이 콘텐트 니고시에이션
* 아래 4개의 헤더를 사용함
  * Accept : MIME 타입을 협상함 / 응답 : Content-Type
  * Accept-Language : 표시 언어 협상 / 응답 : Content-Language 헤더, html 태그
  * Accept-Charset : 문자셋 / 응답 : Content-Type
  * Accept-Encoding : 바디압축 / 응답 : Content-Encoding

#### 파일 종류

* `Accept: image/webp,*/*;q=0.8`
  * webp를 지원하면 webp, 아니면 다른 포맷(우선 순위 0.8) 으로 줄 것을 서버에 요청하는 것
* q는 품질계수 (0~1, 1이면 생략)
* 서버는 요청에서 요구한 형식 중에서 우선순위를 해석하여 가장 일치하는 포맷으로 반환함.
* 서로 일치하는 형식이 없으면 서버가 406 Not Acceptable 오류를 반환함

> 내의견 : 이라고 되어있으나 실제 각종 서버에 엉뚱한 Accept 헤더를 날려보면 그냥 Accept 헤더를 무시하고 서버가 적절히 응답한다. 사용자경험상 더 편하기 때문일듯. 이건 Accept-Language 도 마찬가지. 문자셋과 인코딩도 마찬가지가 아닐까 싶다.

#### 표시 언어 결정

* `Accept-Language: ko-KR,ko;q=0.8,en;q=0.6`
  * ko-KR, ko, en 의 우선순위로 언어 요청
* 품질계수는 Accept와 동일
* 대응되는 헤더로 `Content-Language` 가 있으나 잘 안 쓰임

* HTML의 경우
```html
<!-- 대신 html 태그에 lang이 지정되는 경우는 종종 볼 수 있음 -->
<html lang="ko">
```

#### 문자셋 결정

* `Accept-Charset: utf-8;windows-949;q=0.7;*;q=0.3`
  * utf-8, windows-949, 아니면 그 외의 다른 인코딩으로 응답해줄 것을 요청
* 현대의 모던 브라우저들은 대부분 송신하지 않음. 저자 의견은 브라우저들이 모든 문자셋 인코더를 내장하고 있으므로 미리 네고시에이션을 할 필요가 없어졌기 때문이라 함.
* 컨텐츠의 문자셋은 `Content-Type: text/html; charset=UTF-8` 처럼 MIME타입과 세트로 `Content-Type` 에 실려 통지됨
* 사용할 수 있는 문자셋은 [IANA](https://www.iana.org/assignments/character-sets/character-sets.xhtml)에서 관리함

* HTML의 경우
```html
<!-- RFC 1866 HTML/2.0 스타일 -->
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">

<!-- HTML5 스타일 -->
<meta charset="UTF-8">
```

* 구글 홈페이지의 HTML 시작부분
```html
<!doctype html>
<html itemscope="" itemtype="http://schema.org/WebPage" lang="ko">
  <head>
    <meta charset="UTF-8">
    <!-- 이하 생략 -->
```

#### 압축을 이용한 통신 속도 향상

* 압축의 효과가 큰 리소스(텍스트 등)라면 압축을 통해 통신 비용(시간, 금액 등)을 절감할 수 있다.
* `Accept-Encoding: deflate, gzip`
  * deflate 또는 gzip 으로 압축하여 응답해줄 것을 요청
* 서버는 `Accept-Encoding` 헤더를 보고 지원가능한 인코딩으로 압축하여 응답한다.
*  `Content-Encoding: gzip`
  *  서버는 인코딩에 사용한 알고리즘을 `Content-Encoding` 헤더에 담아 응답한다.
  *  리소스가 인코딩된 경우의 `Content-Length` 값은 압축된 파일의 크기임
* 웹 브라우저에서 사용하는 주요 압축 알고리즘
  * deflate, gzip, br
  * compress, exi, identity (무압축을 선언하는 예약어)
  * sdch : Shared Dictionary Compressing for HTTP. 미리 교환한 사전을 이용한 압축방식. 크롬에서 쓰인다고 함 (IANA 등록 표준X).
    * 이와 같은 공유 사전 방식의 압축은 HTTP/2의 헤더 압축에도 쓰임.
* 이와 같은 Content-Encoding 은 압축을 통해 콘텐츠 크기를 줄이는 방식
* Transfer-Encoding 은 통신 경로를 압축하는 방법인데 그다지 쓰이지 않음.

```bash
# 버전에 따라 다른데 curl/7.68.0 에서는 --compressed 옵션을 사용하면 헤더에
#  Accept-Encoding: deflate, gzip, br
# 을 추가한다.
curl --compressed http://localhost:18888
```

> 내의견 :  이와 반대로 서버가 응답할 때 클라이언트에 Accept-Encoding을 주고 클라이언트에서 서버로 리소스를 보낼 때 압축해서 보내는 방식도 논의되고 있다고함. 21년기준으로 어떨지는 모르겠음. 근데 이건 자체구현하던 서버 기능으로 지원하면 될 일이라 필요에 따라 찾아서 적용하면 될 듯. 예를 들면 json 데이터를 서버로 빈번하게 보내는 사이트에서 유용할듯.

## 쿠키

* 웹사이트의 정보를 브라우저에 저장하는 작은 파일.
* 헤더를 기반으로 구현됨
* HTTP는 stateless 하지만 쿠키를 통해 stateful 처럼 보이게 서비스를 제공할 수 있음 

```
# 서버가 클라이언트에게 쿠키를 저장하도록 헤더에 지정한 경우
Set-Cookie: LAST_ACCESS_DATE=Jul/20/2019
Set-Cookie: LAST_ACCESS_TIME=12:04
```

* 기본적으로 `이름=값` 의 형식
* 클라이언트는 이 값을 저장해두고 다음 요청에는 아래와 같이 전달

```
Cookie: LAST_ACCESS_DATE=Jul/20/2019; LAST_ACCESS_TIME=12:04
```

```bash
# -c : 수신한 쿠키를 지정한 파일에 저장
# -b : 지정한 파일에서 쿠키를 읽어와 전송
# -c, -b를 동시에 사용하면 브라우저처럼 동시에 송수신가능
curl -v -c cookie.txt -b cookie.txt http://localhost:18888/cookie

# -b 는 개별 쿠키 추가에도 사용할 수 있음
curl -v -c cookie.txt -b cookie.txt -b "ABC=MYCOOKIE" http://localhost:18888/cookie
```

#### 쿠키의 잘못된 사용법

* 쿠키는 브라우저 설정에 따라 언제든지 삭제될 수 있고 심지어 아예 저장되지 않을 수도 있음. 따라서 사라지더라도 문제가 없는 정보나 서버로부터 복구 가능한 데이터를 저장하는데 적합함.
* 쿠키의 최대 크기는 4KB 로 정해져 있으므로 주의 필요.
* 쿠키는 항상 통신에 부가되므로 그만큼 전체적인 통신 비용을 높임.
* HTTP 에서는 평문으로 노출되므로 주의. 또한 암호화를 하더라도 사용자가 자유롭게 제어가능하므로 민감한 정보를 저장하는 용도로 사용하면 안 됨

#### 쿠키에 제약을 주다 - 쿠키 옵션

* 기본 양식인 `이름=값` 뒤에 세미콜론을 구분자로 다양한 옵션을 추가할 수 있음.
  * 대소문자 구분 안 함
* 설명에서 알 수 있듯 보안을 위한 옵션이 대부분
  * Expires : 쿠키의 수명설정. `Wed, 26-May-2021 08:11:54 GMT` 형식
  * Max-Age : 초 단위로 지정. 현재 시각에서 지정된 초수를 더한 시간에 쿠키가 만료됨.
  * Domain : 클라이언트에서 쿠키를 전송할 대상 서버. 생략하면 쿠키 발행 서버.
  * Path : 클라이언트에서 쿠키를 전송할 대상 서버의 경로. 생략하면 쿠키를 발행한 서버 경로.
  * Secure : https 일 때만 서버로 쿠키를 전송함.
  * HttpOnly : 자바스크립트 엔진으로부터 쿠키를 숨김. XSS 등의 공격에 대한 방어책.
  * SameSite : RFC에는 없고 크롬에 있는 기능. 같은 오리진(출처)의 도메인에 전송하게 함.

```bash
# ex
Set-Cookie: 1P_JAR=2021-04-26-08; expires=Wed, 26-May-2021 08:11:54 GMT; path=/; domain=.google.com; Secure; SameSite=none
```