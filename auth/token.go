package auth

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt"
)

/*
    EventToken struct defines the event token
    PublisherID -> verifies the publisher
    EventType   -> specifies format of event data
    ExpiresAt   -> the extended time permissions
    Hash        -> the actual event data
*/

type EventToken struct {
    PublisherID string    `json:"publisher_id"`
    EventType   string    `json:"event_type"`
    ExpiresAt   time.Time `json:"expires_at"`
    Hash        string    `json:"hash"`
}

/*
    TokenManager struct manages the verification of the eventSecret and the jwtSecret
    jwtSecret   -> the JWT access token
    eventSecret -> the event access token
*/

type TokenManager struct {
    jwtSecret   []byte
    eventSecret []byte
}

// initializes a new tokenmanager
func NewTokenManager(jwtSecret, eventSecret []byte) *TokenManager {
    t := TokenManager{
        jwtSecret:   jwtSecret,
        eventSecret: eventSecret,
    }
    return &t
}


/*
    createJWT() TokenManager method
    creates a jwt access token that verifies the publisher 
*/
func (tm *TokenManager) CreateJWT(publisherID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
        ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
        IssuedAt:  time.Now().Unix(),
        Issuer:    publisherID,
    })
    
    jwt, err := token.SignedString(tm.jwtSecret)
    return jwt, err
}

/*
    CreateEventToken() TokenManager method
    creates a new event token around the event data
*/
func (tm *TokenManager) CreateEventToken(publisherID, eventType string, payload []byte) (EventToken, error) {
    h := hmac.New(sha256.New, tm.eventSecret)
    h.Write(payload)
    hash := base64.StdEncoding.EncodeToString(h.Sum(nil))
    
    e := EventToken{
        PublisherID: publisherID,
        EventType:   eventType,
        ExpiresAt:   time.Now().Add(24 * time.Hour),
        Hash:        hash,
    }

    return e, nil
}

/*
    VerifyJWT() TokenManager method
    verifies the signature of the access token using the hs256 method
*/
func (tm *TokenManager) VerifyJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return tm.jwtSecret, nil
    })
}

/*
    VerifyEventToken() TokenManager method
    Verifies the expiry date of the token
    Verifies the validity of the hash of the event data
*/
func (tm *TokenManager) VerifyEventToken(token EventToken, payload []byte) error {
    if time.Now().After(token.ExpiresAt) {
        return fmt.Errorf("event token expired")
    }

    h := hmac.New(sha256.New, tm.eventSecret)
    h.Write(payload)
    hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

    if hash != token.Hash {
        return fmt.Errorf("invalid payload hash")
    }

    return nil
}
