package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
)

func main() {
	publicKeyPath := flag.String("public", "public_key.pem", "Путь для сохранения открытого ключа")
	privateKeyPath := flag.String("private", "private_key.pem", "Путь для сохранения приватного ключа")

	flag.Parse()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}

	publicKey := &privateKey.PublicKey
	privateKeyFile, err := os.Create(*privateKeyPath)
	if err != nil {
		fmt.Println("Error creating private key file:", err)
		return
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		fmt.Println("Error encoding private key:", err)
		return
	}

	publicKeyFile, err := os.Create(*publicKeyPath)
	if err != nil {
		fmt.Println("Error creating public key file:", err)
		return
	}
	defer publicKeyFile.Close()

	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Error encoding public key:", err)
		return
	}

	err = pem.Encode(publicKeyFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	if err != nil {
		fmt.Println("Error encoding public key:", err)
		return
	}

	fmt.Println("Ключи успешно сгенерированы и сохранены в файлы")
	fmt.Printf("- Private key: %s\n", *privateKeyPath)
	fmt.Printf("- Public key: %s\n", *publicKeyPath)
}
