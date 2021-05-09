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

## 인증과 세션

* Basic 인증 : 가장 간단한 인증방식
  * 유저명, 패스워드를 BASE64로 인코딩하여 전송. 평문으로 감청될 경우 손쉽게 인증정보가 탈취됨.
  * base64(유저명 + ":" + 패스워드)
```bash
# 기본인증방식이 basic이므로 --basic 생략가능 
curl --basic -u user1:pw123! http://localhost:1888

# 서버단 헤더
# ...
# Authorization: Basic dXNlcjE6cHcxMjMh
```

* Digest 인증 : 해시함수를 이용함 
  * id, pw, uri, 요청마다 달라지는 nonce 값 등으로 계산한 해시값을 사용하므로 노출되도 복호화가 어려움
```bash
curl --digest -u user1:pw123! http://localhost:18888
```

* Basic, Digest 인증 모두 일반적으로는 쓰이지 않음
  * 특정 경로 아래를 보여주지 않는 방식으로만 인증 가능. 이 때문에 톱페이지부터 인증이 되야 표시가 가능함.
  * 요청할 때마다 인증정보를 보내야함
  * 로그인 화면의 사용자화가 불가함
  * 명시적인 로그오프가 불가함
  * 로그인한 단말을 식별할 수 없음. 동시 로그인 불가 등의 기능을 구현할 수 없음.
* 이 때문에 흔히 알고 있듯 form을 이용해 ID/PW (+2차 인증 등) 로 1회 로그인하고 서버는 세션 토큰 `(JSESSIONID 등)` 을 발행함. 세션 토큰은 쿠키로 클라이언트에 발급되고 클라이언트는 이후의 요청에 쿠키를 전송함으로서 로그인 상태가 유지됨. 사이트간 요청 위조 (Cross-Site Request Forgery, CSRF) 대책으로 랜덤 키를 같이 보내기도 함.

> 책에서는 서명된 쿠키를 이용한 세션 데이터 저장도 소개했음. 통신 속도가 빨라진 만큼 쿠키의 데이터양 증가는 큰 문제가 아니니 쿠키에 세션 정보 전체를 저장하고 서버에는 세션 스토리지를 두지 않는다는 아이디어임. 쿠키값 위조가 걱정되는데 이는 쿠키값을 서버만 가지고 있는 공개키, 비밀키로 암호화하는 방식으로 해결한다고 함. 모든 키를 서버가 가지고 있으므로 클라이언트가 변조해봐야 서버는 이를 즉시 알아차릴 수 있을 것임. 또한 세션 스토리지 암호화 방식만 공통화해두면 마이크로서비스에서도 사용할 수 있을만큼 범용성도 좋음. 루비 온 레일즈나 장고에서 지원한다고 함. 장점은 서버측에 세션 스토리지를 두지 않아도 된다는 점이고 단점은 네트워크 트래픽의 증가와 같은 사용자가 다른 클라이언트로 접속했을 때 데이터 공유가 어렵다는 점.

> 내의견) 서버에서 세션 스토리지를 두려면 HA 구성이든 뭐든 돈이 드는데 이 방식을 도입한 비용절감 효과가 네트워크 트래픽 증가에 따른 비용보다 크다면 나름 도입할만한 방법인 것 같다. 근데 또 떠오른 단점은 현재 접속한 세션수를 세는게 어렵겠다.

## 프록시

* HTTP/1.0 에서는 프록시와 게이트웨이를 다음과 같이 정의한다.
  * 프록시 : 통신 내용을 이해한다. 필요에 따라서 콘텐츠를 수정하거나 서버 대신 응답한다.
    * ex) 컨텐츠 캐싱, 외부 공격으로부터 네트워크를 보호하는 방화벽 역할, 저속 통신 회선용으로 데이터 압축하는 필터나 콘텐츠 필터링 등
  * 게이트웨이 : 통신 내용을 그대로 전송한다. 내용의 수정도 불허한다. 클라이언트에서는 중간에 존재하는 것을 알아채서는 안 된다.
* HTTPS 통신의 프록시 지원은 HTTP/1.1에서 추가된 CONNECT 메서드를 이용함

* 프록시 구조는 단순해서 GET 등의 메서드 다음에 오는 경로명 형식만 바꿈.

```bash
# 원래 요청
GET /helloworld
Host: localhost:18888

# 프록시를 설정하면 요청 경로에 스키마가 추가되므로 URL 형식이 됨
GET http://example.com/helloworld
Host: localhost:18888
```

* 프록시 서버로 보내는 경우. 프록시는 이를 받아 중계할 곳으로 요청을 리디렉트하고 결과를 클라이언트에 반환함
* 프록시 서버에서 인증을 사용하는 경우 프록시 서버는 `Proxy-Authenticate` 헤더로 인증이 필요함을 클라이언트에게 알리고 클라이언트는 인증정보를 `Proxy-Authorization` 에 담아 전송함.
* 중계되는 프록시는 중간의 호스트 IP 주소를 특정 헤더에 기록함. 옛날부터 쓰던 비표준은 `X-Forwarded-For` 헤더이고 표준은 RFC 7239에 추가된 `Forwarded` 헤더임. 다만 남길지 말지는 프록시 서버 마음이므로 그대로 믿으면 안 됨.
```
X-Forwarded-For: client, proxy1, proxy2
```

```bash
# -x/--proxy : 프록시 서버를 지정함
# -U/--proxy-user 프록시 서버 인증정보
# --proxy-basic, --proxy-digest 등으로 프록시 인증 방식 변경 가능
curl -x http://localhost:18888 -U user:pass http://google.com/

':
GET http://google.com/ HTTP/1.1
Accept: */*
Proxy-Authorization: Basic dXNlcjpwYXNz
Proxy-Connection: Keep-Alive
User-Agent: curl/7.68.0
'
```

## 캐시

#### 갱신일자에 따른 캐시

* HTTP/1.0 시절에는 정적 콘텐츠 위주이므로 콘텐츠가 갱신됐는지만 비교하면 충분했음

```bash
# 웹서버는 아래 헤더를 응답에 포함시킴.
# 날짜는 RFC 1123 으로 기술되며, 타임존에는 GMT를 설정
Last-Modified: Wed, 08 Jun 2020 15:23:45 GMT

# 웹 브라우저는 캐시된 URL을 다시 읽을 때 서버가 반환한 Last-Modified 값을 그대로 아래 헤더에 담아 요청함
If-Modified-Since: Wed, 08 Jun 2020 15:23:45 GMT

# 콘텐츠가 변경됐으면 200 OK 응답과 함께 콘텐츠를 응답 바디에 실어서 보냄
# 변경되지 않았으면 304 Not Modified 를 반환하고 바디를 응답에 포함하지 않음
```

![cache2.8.1](./imgs/realworldhttp/2.8.1.cache.png)

#### Expires

* `Last-Modified`, `If-Modified-Since` 을 이용한 방식은 어쨋든 서버로 요청이 발생함
* Expires 는 콘텐츠의 유효기간을 정함으로써 지정한 기간 내에는 강제로 캐시를 이용하도록 하고 서버로 요청을 아예 전송하지 않음
* 엄밀히는 Expires 에 지정된 시간은 서버에 접속을 할지 말지 판단할 때만 사용함
* 또한 브라우저의 "뒤로 가기 버튼" 등으로 방문 이력을 조작하는 경우는 기한이 지난 오래된 콘텐츠가 그대로 이용될 수도 있음
* 지정된 시간까지는 서버에 아예 요청을 보내지 않으므로 주의해서 사용해야함. RFC 2068 에서는 변경할 일이 없더라도 최대 1년의 캐시 수명을 설정하자고 가이드하고 있음.

```bash
# 지정된 시간 이내에 재요청이 발생하면 서버로 요청하지 않고 가지고 있는 캐시를 그대로 사용
# 이후에 재요청이 발생하면 서버에 재요청
Expires: Wed, 08 Jun 2020 15:23:45 GMT
```

> 내의견 : Expires 를 사용하더라도 해당 콘텐츠를 로드하는 페이지 자체는 캐싱을 하지 않도록 설정해두고 그 페이지에서 Expires 헤더를 사용한 콘텐츠의 url 뒤에 ?v=20210508 과 같은 파라미터를 붙이면 URL이 달라지므로 다시 서버로 요청이 들어간다. 

* 아래 그림에서 "서버에 접속"은 전술한 `Last-Modified` 헤더를 이용한 캐시 로직이 들어감

![cache2.8.2](./imgs/realworldhttp/2.8.2.cache.png)

#### Pragma: no-cache

* Pragma는 지시를 포함한 요청 헤더가 들어가는 헤더
* Pragma 헤더에 포함할 수 있는 페이로드로 유일하게 HTTP 사양으로 정의된 것이 no-cache
* no-cache는 "요청한 콘텐츠가 이미 저장돼 있어도, 원래 서버(오리진 서버)에서 가져오라" 고 프록시 서버에 지시하는 것. HTTP/1.1 에서 Cache-Control로 통합됐으나 하휘 호환성 유지를 위해 남아있음
* 프록시가 요청한 대로 처리하리라는 보장은 없음. 중간에서 프록시가 하나라도 no-cache를 무시하면 기대한 대로 동작하지 않음.
* HTTP/2 부터는 프록시가 통신 내용을 감시할 수 없고 중계만할 수 있으므로 프록시의 캐시를 외부에서 관리하는 의미는 이제 없다고도 말할 수 있음 
* 이런 이유로 별로 사용되지 않음

#### ETag

* 날짜와 시간을 이용한 캐시 비교만으로는 해결할 수 없는 상황도 있음
* 동적으로 바뀌는 요소가 늘어날수록 날짜를 근거로 캐시의 유효성을 판단하기 어려움
  * 사용자마자 화면 구성이 동적으로 달라지는 사이트 등
* RFC 2068의 HTTP/1.1에서 ETag(Entity Tag) 가 추가됨. 이 값은 순차적인 갱신 일시가 아니라 파일의 해시 값으로 비교함
* 서버는 응답에 `ETag` 헤더를 부여함. 두 번째 이후 요청시 클라이언트는 `If-None-Match` 헤더에 캐시에 있던 ETag 값을 추가해 요청함. 서버는 파일의 ETag 값과 비교해 같으면 304 Not Modified 로 응답함. 즉, 비교값만 달라졌을 뿐 `Last-Modified` 를 사용한 방식과 동일함.
* 대신 ETag는 서버가 자유롭게 결정해서 반환할 수 있음.
  * 콘텐츠 파일의 해시값, "갱신일시-파일크기" 형식의 스트링 등
  * 과거에는 inode 값을 사용한 적도 있었는데 서버를 여러 대 병렬시킨 경우 같은 콘텐츠인데도 id가 달라 ETag도 달라지므로 현재는 사용하지 않는 방식임

![cache2.8.3](./imgs/realworldhttp/2.8.3.cache.png)

#### Cache-Control (1)

* ETag와 같이 HTTP/1.1에 추가됨
* 더 유연한 캐시 제어를 지시할 수 있음
* Expires보다 우선해서 처리됨

---

* 서버에서 응답을 보낼 때는 아래와 같은 키를 사용할 수 있음
  * public : 같은 컴퓨터를 사용하는 복수의 사용자간 캐시 재사용 허가
  * private : 같은 컴퓨터를 사용하는 다른 사용자간 캐시 재사용하지 않음. 같은 URL에서 사용자마자 다른 컨텐츠가 돌아오는 경우 이용
  * max-age=n : 캐시의 신선도를 초단위로 설정. 86400이면 하루동안 캐시가 유효하고 서버에 문의하지 않고 캐시를 이용함. Expires의 역할을 한다고 볼 수 있음. 그 이후는 서버에 문의한 뒤 304 Not Modified 가 반환됐을 때만 캐시를 이용함
  * s-maxage=n : max-age와 동일하나 공유 캐시에 대한 설정값
  * no-cache : 캐시가 유효한지 매번 문의함. max-age=0가 거의 동일함.
  * no-store : 캐시하지 않음
* no-cache는 캐시하지 않는다는 말이 아니고 항상 서버로 문의하여 갱신 일자와 ETag를 사용한 캐시 정책을 사용하겠다는 의미임. 캐시하지 않는 것은 no-store 임.
* 콤마로 구분해 복수 지정이 가능함. 보통 아래와 같이 조합
  * private, public 중 하나. 혹은 생략 (기본값 private)
  * max-age, s-maxage, no-cache, no-store 중 하나

```bash
#구글 메인화면의 js 중 하나의 캐시 관련 헤더
cache-control: public, max-age=31536000
expires: Sun, 08 May 2022 01:44:52 GMT
last-modified: Fri, 07 May 2021 20:29:55 GMT
vary: Accept-Encoding
```

---

* 아래를 포함한 cache 설명 이미지는 설명을 위한 대략적인 설명이며 항상 맞다는 보증은 없음. 예를 들어 모순된 설정을 동시에 할 경우 (no-cache와 max-age) 의 우선순위까지는 RFC에 적혀있지 않음.

![cache2.8.4](./imgs/realworldhttp/2.8.4.cache.png)

#### Cache-Control (2)

* `Cache-Control` 헤더는 클라이언트 요청헤더로 쓰일 때는 `Pragma: no-cache` 처럼 프록시에 다양한 지시를 할 수 있음
    * no-cache: `Pragma: no-cache`와 동일
    * no-store : 응답의 no-store와 같고, 프록시 서버에 캐시를 삭제하도록 요청
    * max-age : 프록시에서 저장된 캐시가 최초로 저장되고 나서 지정 시간 이상 캐시는 사용하지 않도록 프록시에 요청
    * max-stale, min-fresh, no-transform, only-if-cached 등

---

* 한편 응답 헤더에서 서버도 프록시에게 캐시 컨트롤 지시를 내릴 수 있음. 물론 전항의 서버에서 클라이언트로 보내는 지시에서 소개한 명령은 모두 프록시에도 유효하고 아래 키들은 프록시 서버를 위한 전용 명령들임
  * no-transform : 프록시가 콘텐츠를 변경하는 것을 제어함
  * must-revalidate, proxy-revalidate 등

#### Vary

* 같은 URL 이라도 클라이언트에 따라 반환 결과가 다름을 나타내는 헤더
* 예를 들어 모바일 브라우저, 데스크탑 브라우저에 따라 표시가 달라질 수 있음
* 이처럼 표시가 바뀌는 이유에 해당하는 헤더명을 Vary 헤더에 나열함으로써 잘못된 콘텐츠의 캐시로 사용되지 않게 함

```bash
# User-Agent와 Accept-Language에 따라 콘텐츠가 달라질 수 있음을 명시
Vary: User-Agent, Accept-Language
```

* 참고로 User-Agent 는 관례일 뿐 정규화된 정보가 아니므로 판정이 틀릴 수도 있음.  2017년 구글 가이드라인에서는 같은 콘텐츠 (HTML, CSS, JS) 를 모든 브라우저에 배포하고, 브라우저가 필요한 설정을 선택하는 반응형 웹 디자인을 권장함.

## 리퍼러 (Referer)

* 사용자가 어느 경로로 서버에 도달했는지 파악할 수 있도록 클라이언트가 서버에 보내는 헤더.
* 원래 스펠링은 referrer 인데 RFC 1945 제안 당시의 오자가 남은 것이라고함.
* GET 파라미터에 개인정보 등 민감한 정보가 있을 경우 리퍼러를 통해 타사이트로 그대로 유출될 수 있으므로 파라미터에 민감한 정보가 노출되지 않도록 해야함.
* 리퍼러의 용도 예제
  * 어떤 사이트로부터 우리 서비스로 들어온건지 파악할 때 사용
  * 이미지가 타사이트에 직접 링크되는 것을 막을 때 사용
  * CSRF 방어 목적으로도 사용했으나 브라우저에서 리퍼러를 전송하지 않도록 설정할 수도 있으므로 권장하지 않음

```bash
# 구글에서 검색 후 결과 페이지에서 특정 사이트로 이동하면 해당 요청에는 아래와 같은 Referer 가 실린다.
Referer: https://www.google.com/

# 구글 검색결과 페이지에는 아래의 referrer 정책이 정의되어있으므로 도메인 이름만 전송되었음을 알 수 있다.
# ... <meta content="origin" name="referrer"> ...
```

* 스키마 조합과 리퍼러의 유무

|액세스 출발지|액세스 목적지|리퍼러를 전송하는가?|
|---|---|---|
|HTTPS|HTTPS|한다|
|HTTPS|HTTP|하지 않는다|
|HTTP|HTTPS|한다|
|HTTPS|HTTPS|한다|

* 대부분의 브라우저는 이 규칙을 준수함. 단, 이 규칙을 엄밀히 적용할 경우 서비스간 연계에 차질이 생기기도 해서 IE, 낮은 버전의 안드로이드 브라우저 등 준수하지 않는 브라우저도 존재함.

---

* 리퍼러 정책 설정 방법 (Referrer 로 오자가 수정됐으므로 주의)
  * Referrer-Policy 헤더
  * Content-Security-Policy 헤더
  * `<meta name="referrer" content="설정값">` 태그
  * `<a>` 태그 등 몇 가지 요소의 referrerpolicy 속성 및 `rel="noreferrer" 속성
* 정책 설정 값 예제
  * no-referrer : 전혀 보내지 않음
  * no-referrer-when-downgrade : 현재 기본 동작처럼 HTTPS -> HTTP 일 때는 전송하지 않음
  * same-origin : 동일 도메인 내의 링크에 대해서만 전송
  * origin : 도메인 이름만 전송
  * strict-origin : origin과 같지만 HTTPS -> HTTP 일 때는 전송하지 않음
  * origin-when-crossorigin : 같은 도메인 내에서는 완전한 리퍼러를, 다른 도메인에는 도메인 이름만 전송
  * strict-origin-when-crossorigin : origin-when-crossorigin :과 같지만 HTTPS -> HTTP 일 때는 전송하지 않음
  * unsafe-url : 항상 전송

---

* Content-Security-Policy 예제
* CSP 헤더는 다양한 보안 설정을 한번에 변경할 수 있는 헤더임. 10장에서 자세히 설명함.

```bash
Content-Security-Policy: referrer origin
```

## 검색 엔진용 컨텐츠 접근 제어

* 크롤러(로봇, 봇, 스파이더 등)의 접근을 제어하는 방법을 주로 다음 두 가지가 쓰임
  * robots.txt
  * sitemap
* 미국 재판에서는 robots.txt 가 법적 효력을 가진다는 판례가 여러개 있다고 함. 자세한건 robots.txt 웹사이트 참고

```bash
# robots.txt 예제
User-agent: *
Disallow: /service/

# 메타태그로도 설정 가능
# 아래 예제는 검색엔진이 인덱스하는 것을 거부함을 의미
<meta name="robots" content="noindex">
```
---

* sitemap은 웹사이트에 포함된 페이지 목록과 메타데이터를 제공하는 XML 파일
  * https://www.sitemaps.org
* 주로 검색 엔진에 정보를 제공하는 용도임

## 마치며

* HTTP는 효율적으로 계층화되어 있음. 통신의 데이터 상자 부분은 변하지 않으므로, 규격에서 제안된 새로운 기능이 구현되지 않아도 호환성을 유지하기 쉽도록 되어 있음. 또한 압축 방식 선택 등 브라우저가 규격화되지 않은 방식을 새로 지원해도 가능하다면 사용할 수 있음. 토대가 되는 문법(신택스)과 그 문법을 바탕으로 한 헤더의 의미 해석(시맨틱스)이 분리되어 있으므로 상위 호환성과 하위 호환성이 모두 유지됨.

# Ch03 Go 언어를 이용한 HTTP/1.0 클라이언트 구현

* 이 장은 curl 로 해본 예제들을 Go 언어로 구현해보는 파트이므로 별도 정리는 생략함

# Ch04 HTTP/1.1의 신택스: 고속화와 안정성을 추구한 확장

