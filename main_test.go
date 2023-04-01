package main

import (
	"testing"

	"gopkg.in/h2non/gock.v1"
)

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4vwlqGX9ZmiJk6BGs/ZCTxhJclSYJMrKnDvQBQjUWg9Cwg2n
LJWW+BUnQeF7P4PbyVx7BS/mJ0POjeP3Rrzjp59C9wAbM8VhvnYNm3xUCQJFLWut
uejisZ19UrV/ZXVb4IDec+0tBd9WVFzqWmRCbq3cW9gIVJpAdqJeFcFY71L9knF9
fCuPekIcQHd0Oic4USIF/kJXyR4ctFKBQir4ctwCHlALl3AxgeKFr7me88X9Ob5v
HydQqyWl6oBIHzyIzLfveP04NO5WVbda4LrbjdASv96V3iu7U5BzED/zGmpcaZzD
6aa53oNVdurrZ/kIKz+GVymMFKCvM37GAQBM3wIDAQABAoIBAQCGdIBGGWwaTpA4
L3fSQGyU97kCDZQ2Lx4Hn/KgGNPZKTMNShMd+Np9x+ICR3O/cvctdye0MeRum97t
8/zVHSzpbRC4yYpTh3dX4Aw9b09EKuEZf7Bf8NDgD39eD/8P9Y3gFdYv622BDgPQ
Y126/6rObxSaHwUIQHsxCwsabfalhHCmVjdSJ+Upam7z7ZZohS9qWxabrEC3urA3
Dqv+Rnfg5qc2ZrgRTlHXnBnTAAgMlrSMbMgAn2GhS3KUIBfacEKWblLJR2FfcXcq
zeUwWaWPIWIkDFpFrF68HEuKpjapySUlbaK2cgVSn6l6qA4gUAEA0hYGzNYOmhR2
HZ65Cbr5AoGBAPx5MoFUFgLtu37sVZU8OtezQWVPC3R2cbcNbq1kTw/GyHzHHlWe
Pf52orY79MByQFrEbEjuWe+Xt64SsM6+quMMywaDCj27pI14HrWbDPplxHphLIqY
wilPWkA72E5nA04OfSGIIfgMW0h3OHiOduyMOlXixjXgaI7KvKZ0GH2rAoGBAOYn
zqpUyTXj6RX6EcdRQKJTdpjYZawl4b1CgS6eYpq8y8sMoCFlQWX/xNKQaMj6U2yW
iNi+Yy6LnjafzM9lstflNDhW+/L1FdtGerNbK92TwGJ29aJaxG52icBigi8C1u/E
egW7LWjdf7WEDJNHQDiyaDl5+yB6DX5bnyyVSLGdAoGAVi2AYciz4rgHAeHlrJTs
eOgE8HG0tUIgupzpJGJS4k217XGCFzN2cb9I9u8sMexNry3Q0GwbYr7kwZQ7qbZH
WkzpmAVun3fHSUqxIMgV+/p0wFke/Qf7bmJZqgdDZC+hXylu6N0wyxxcpDWdnvjx
+vg6iUpo4ccBqYvmLOL/4RUCgYAvJX5nVAD3wh0wPE7CBrn3xqMnwkRplET+0Q3H
b/iA/CW/DXIMBUL1UwSNobllWioWt2uHAtEsartZMzjwT0Poh/I/jEoGRgBZL8HY
1ddRh3/Ea9v7ix5sBmpHd6Z1XN6MtTHN1L8DmUQc+dTdop3cP2esRnmT+IylEr2z
k00V3QKBgFLthEgAYGyHh0W3uzG3aDH9jcR/IwrAXiOXKjZtSf161J3vY7aehETb
ZQdFzw/6wDAcr1sOUf4FLKE0W5lRytPsF5Qxl5mgnbVFuNHujUb7IY7dk1St7gde
TXAue3jmAcAqjcVjR0j7iRhmyEAMf5ImHFhYaVhp83bXiu893aV4
-----END RSA PRIVATE KEY-----`)

var token = "ghs_16C7e42F292c6912E7710c838347Ae178B4a"

func Test_fetchTokenReturnsSuccess(t *testing.T) {
	defer gock.Off()
	gock.New("https://api.github.com/app/installations/1/access_tokens").
		MatchHeader("Accept", "application/vnd.github.v3+json").
		MatchHeader("Authorization", "^Bearer .*$").
		Reply(201).
		JSON(map[string]string{"token": token})

	value, _ := FetchTokenFromAPI("1", "2")

	if value != token {
		t.Errorf("Expected %s, got %s", token, value)
	}
}

func Test_failsWhenPrivateKeyIsInvalid(t *testing.T) {
	_, err := FetchTokenFromAPI("1", "2")

	if err == nil {
		t.Errorf("Expected failure, program ran through.")
	}
}

func Test_failsWhenApiReturnsError(t *testing.T) {
	type args struct {
		jwtToken      string
		applicationID string
	}

	tests := []struct {
		code    int
		args    args
		want    string
		wantErr bool
	}{
		{code: 401, args: args{jwtToken: "1", applicationID: "2"}, wantErr: true},
		{code: 404, args: args{jwtToken: "1", applicationID: "2"}, wantErr: true},
		{code: 500, args: args{jwtToken: "1", applicationID: "2"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run("testing", func(t *testing.T) {
			gock.New("https://api.github.com/app/installations/1/access_tokens").
				MatchHeader("Accept", "application/vnd.github.v3+json").
				MatchHeader("Authorization", "^Bearer .*$").
				Reply(tt.code)

			_, err := FetchTokenFromAPI(tt.args.jwtToken, tt.args.applicationID)

			if (err != nil) != tt.wantErr {
				t.Errorf("fetchInstallationToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_generateJWTSucess(t *testing.T) {
	_, err := generateJWT("1", privateKey)

	if err != nil {
		t.Errorf("Failed with: %s:", err)
	}
}
