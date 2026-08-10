package tests

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true

// Token — учебный JWT для проверки защищённых API-запросов.
// Он подписан паролем TODO_PASSWORD=87654321 и намеренно не имеет срока действия,
// чтобы CI мог воспроизводимо запускать тесты. В production пароль и токены нельзя
// хранить в исходном коде.
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwdXJwb3NlIjoiZWR1Y2F0aW9uYWwtdGVzdHMifQ.-PEl5HPtS9Iu9sSgKsQrpE_7Aqtiv7KTEtHozmxU-CE`
