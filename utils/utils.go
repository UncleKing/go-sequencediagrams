package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func IntersectSlices(l1, l2 []string) []string {
	var l3 []string
	fm := make(map[string]bool)
	for _, item := range l1 {
		fm[item] = true
	}

	for _, item := range l2 {
		if _, ok := fm[item]; ok {
			l3 = append(l3, item)
		}
	}
	return l3
}

func GetTenantID(c *gin.Context) (string, error) {
	t := c.GetHeader("tenantID")
	if len(t) == 0 {
		return "", fmt.Errorf("Invalid Tenant ID or tenant not specified")
	}
	return t, nil
}

//SwapMax ensures that first number is lesser than second.
func SwapMax(n1, n2 float64) (float64, float64) {
	if n1 > n2 {
		return n2, n1
	} else {
		return n1, n2
	}
}
