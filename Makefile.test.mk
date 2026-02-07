.PHONY: test

HOST := http://127.0.0.1:8000

test: health test_

health:
	curl -X GET "$(HOST)/health"

test_:
	curl -X GET "$(HOST)/test"
