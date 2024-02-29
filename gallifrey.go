package main

import (
    "bufio"
    "crypto/ed25519"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "os"
)

func GenerateEd25519Key() (ed25519.PublicKey, ed25519.PrivateKey) {
    publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error generating Ed25519 key: %v\n", err)
        os.Exit(1)
    }
    return publicKey, privateKey
}

func SignMessage(privateKey ed25519.PrivateKey, message []byte) string {
    signature := ed25519.Sign(privateKey, message)
    return base64.StdEncoding.EncodeToString(signature)
}

func VerifySignature(publicKey ed25519.PublicKey, message []byte, signature string) bool {
    sigBytes, err := base64.StdEncoding.DecodeString(signature)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error decoding signature: %v\n", err)
        return false
    }
    return ed25519.Verify(publicKey, message, sigBytes)
}

func SaveKeyToFile(keyBytes []byte, filename string) {
    encodedKey := base64.StdEncoding.EncodeToString(keyBytes)
    err := os.WriteFile(filename, []byte(encodedKey), 0600)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to save key to %s: %v\n", filename, err)
        os.Exit(1)
    }
}

func LoadKeyFromFile(filename string) []byte {
    encodedKey, err := os.ReadFile(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load key from %s: %v\n", filename, err)
        os.Exit(1)
    }
    keyBytes, err := base64.StdEncoding.DecodeString(string(encodedKey))
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to decode key from %s: %v\n", filename, err)
        os.Exit(1)
    }
    return keyBytes
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: gallifrey <mode> [arguments]")
        fmt.Println("Modes:")
        fmt.Println("\tsign")
        fmt.Println("\tverify <publicKey> <signature>")
        os.Exit(1)
    }

    mode := os.Args[1]
    reader := bufio.NewReader(os.Stdin)

    switch mode {
    case "sign":
        message, _ := reader.ReadBytes('\n')

        var privateKey ed25519.PrivateKey
        var publicKey ed25519.PublicKey

        if _, err := os.Stat("private_key.pem"); os.IsNotExist(err) {
            publicKey, privateKey = GenerateEd25519Key()
            SaveKeyToFile(privateKey, "private_key.pem")
            SaveKeyToFile(publicKey, "public_key.pem")
        } else {
            privateKeyBytes := LoadKeyFromFile("private_key.pem")
            privateKey = ed25519.PrivateKey(privateKeyBytes)
            publicKeyBytes := LoadKeyFromFile("public_key.pem")
            publicKey = ed25519.PublicKey(publicKeyBytes)
        }

        signature := SignMessage(privateKey, message)
        fmt.Printf("Signature: %s\n", signature)
        fmt.Printf("Public Key: %s\n", base64.StdEncoding.EncodeToString(publicKey))

    case "verify":
        if len(os.Args) != 4 {
            fmt.Println("Usage for verify: signVerify verify <publicKey> <signature>")
            os.Exit(1)
        }

        message, _ := reader.ReadBytes('\n')
        publicKeyString := os.Args[2]
        signature := os.Args[3]

        publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyString)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error decoding public key: %v\n", err)
            os.Exit(1)
        }

        publicKey := ed25519.PublicKey(publicKeyBytes)
        isValid := VerifySignature(publicKey, message, signature)
        if isValid {
            fmt.Println("\033[42;30;1m### SIGNATURE OK ###\033[m")
        } else {
            fmt.Println("\033[41;30;1m### SIGNATURE INVALID ###\033[m")
        }

    default:
        fmt.Println("Invalid mode. Use 'sign' or 'verify'.")
        os.Exit(1)
    }
}
