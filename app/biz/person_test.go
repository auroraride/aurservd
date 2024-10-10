// Copyright (C) liasica. 2024-present.
//
// Created at 2024-03-02
// Based on aurservd by liasica, magicrolan@qq.com.

package biz

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/auroraride/adapter/rpc/pb"
	"github.com/golang-module/carbon/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/auroraride/aurservd/internal/ar"
)

func TestPersonIdentity(t *testing.T) {
	// 初始化实名认证rsa客户端
	ar.LoadRsa()

	identity := &pb.PersonIdentity{
		IdCardNumber: "111222198612123333",
		Name:         "张三",
		OcrResult: &pb.PersonIdentityOcrResult{
			Name:            "张三",
			Sex:             "男",
			Nation:          "汉",
			Birth:           "19861212",
			Address:         "北京市天安门",
			IdCardNumber:    "111222198612123333",
			ValidStartDate:  "20191112",
			ValidExpireDate: "20391111",
			Authority:       "北京市东城区分局",
			PortraitCrop:    "AAA",
			NationalCrop:    "BBB",
		},
	}
	encoded := encryptPersonIdentity(identity)
	decoded, err := decryptPersonIdentity(encoded)
	require.NoError(t, err)
	require.True(t, proto.Equal(identity, decoded))
	expireDate := carbon.Parse(identity.OcrResult.ValidExpireDate).StdTime()
	if expireDate.Before(time.Now()) {
		t.Log("证件已过期")
		return
	}
	t.Log("证件未过期")
}

func TestRSAReadPrivateKeyFromFile(t *testing.T) {
	privFile, err := os.Open("person_private_key.pem")
	if err != nil {
		panic(err)
	}
	defer func(privFile *os.File) {
		_ = privFile.Close()
	}(privFile)

	privBytes, err := io.ReadAll(privFile)
	if err != nil {
		panic(err)
	}

	privBlock, _ := pem.Decode(privBytes)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		panic("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		panic(err)
	}

	cipherText, _ := base64.StdEncoding.DecodeString("XT3nMlPDqpA/jWG7d+v9g8WMkPDSdTmPFDgCaIo9fIRIubhmGwuimc5GPKuxsCdxZFFjWvR22KWho5WH8O2AGBsOGEpffLlAjKXbRxi7OSg65JW+w2RHX/wISeEaAyuZqKpE1MNIw0jV0pXzSrNMoQDZkF+ssstg9u/iXDxT6oe5DUMy33UwgxWctc9dw6638eiBskSDm2w6v3yx/OeRRZ9+h11z335cwlUWQuAi6QqxxWTwDdQxryRXK++on2y645lOYpmbRK1zNU4M8IpPJqn5uB5VgqYO4++mreRJsCIU1g7avb6tUbjSo+wnSDHICGHmag01Claq46FCyhNggxnB4tCcnWUzAHM+LgR58UIRhiuHmDrGQMBG+NlAe+Xmo1dgltJ+6cES9A0wvq0vsyr/Tbi8tGgTxuiDzY0Ek72ad6glT0Cc5ZuUgYoAhOthySeJVNoTlT2537GCzPeqtvXZlU6D4AAfCNZEIYjMULASjZVI3JpcQ6tnjrA1grJJWp0Qjby6FMCI6uYeT1XsnxVYJos0xUJKepB5/GVFX922vuBim2iRIsFSkM4l4HCzOT4RvtHZBQHnGdMrE95T55bNFUUkx8HeLHdRqo8Pe9oV3zdqHX9J0gKkBX968C9piJozeknOgJ/pnNE1xpAhuZFcTw5GOCE/N8+ZK6S1Wd8=")
	plainText, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, privateKey, cipherText, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(plainText))
}

func TestProto(t *testing.T) {
	identity := &pb.PersonIdentityOcrResult{
		Name:            "张三",
		Sex:             "男",
		Nation:          "汉",
		Birth:           "19861212",
		Address:         "北京市天安门",
		IdCardNumber:    "111222198612123333",
		ValidStartDate:  "20191112",
		ValidExpireDate: "20391111",
		Authority:       "北京市东城区分局",
		PortraitCrop:    "AAA",
		NationalCrop:    "BBB",
	}
	src, _ := proto.Marshal(identity)
	fmt.Println(base64.StdEncoding.EncodeToString(src))

	err := proto.Unmarshal(src, identity)
	require.NoError(t, err)
}

func TestIDCardNumber(t *testing.T) {
	// 获取生日
	birth := "111015200611051101"[6:14]
	birthday := carbon.Parse(birth).StdTime().AddDate(18, 0, 0)

	// 未年满18岁认证标记为失败
	require.False(t, birthday.After(time.Now()))
}
