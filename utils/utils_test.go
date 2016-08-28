package utils_test

import (
    "encoding/base64"

    utils "github.com/joaodias/hugitoApp/utils"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
    Describe("Check if random generated lengths are right", func() {
        Context("When generating bytes", func() {
            randomByteArray, err := utils.GenerateRandomBytes(24)
            actualLength := len(randomByteArray)
            if err != nil {
                Fail("GenerateRandomBytes returned an error")
            }
            It("should be equal to 24", func() {
                Expect(actualLength).To(Equal(24))
            })
        })

        Context("When generating strings", func() {
            randomString, err := utils.GenerateRandomString(24)
            randomByteArray, _ := base64.URLEncoding.DecodeString(randomString)
            actualLength := len(randomByteArray)
            if err != nil {
                Fail("GenerateRandomString returned an error")
            }
            It("should be equal to 24", func() {
                Expect(actualLength).To(Equal(24))
            })
        })
    })

    Describe("Checking if strings do match", func() {
        Context("With equal strings", func() {
            baseString := "abc"
            stringToCompare := "abc"
            It("should be true", func() {
                Expect(utils.AreStringsEqual(baseString, stringToCompare)).To(Equal(true))
            })
        })
    })

    Describe("Checking if strings do not match", func() {
        Context("With different strings", func() {
            baseString := "abc"
            stringToCompare := "xyz"
            It("should be false", func() {
                Expect(utils.AreStringsEqual(baseString, stringToCompare)).To(Equal(false))
            })
        })
    })
})
