afreeca랑 warudo 연동 예제

별풍선만 읽음. 채널포인트는 나중에.

1. `afreeca-warudo.exe` 와 `.env` 를 다운로드받는다.
2. `.env`에 `BJ_ID`를 수정한다.
3. `afreeca-warudo.exe`를 커맨드라인에서 실행한다.
    - 커맨드라인에서 실행해야 에러가 날 경우에 확인할 수 있습니다.
    - 그대로 실행해도 되지만 에러가 날 경우 그냥 꺼져버림.

이 때 `UDP localhost:19190/osc/balloon` 으로 `Arg1` 로는 보낸 닉네임, `Arg2` 로는 별풍선 갯수를 전달함.

또는

직접 빌드한다.