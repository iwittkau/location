.PHONY: gow-run
gow-run:
	gow run cmd/location/main.go -debug -secret='$(LOCATION_APP_HASH)'