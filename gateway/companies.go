package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"xm/companies"
	"xm/proto/gen"
)

type AddCompanyRequest struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	EmployeesNum uint32                `json:"employees_num"`
	Registered   bool                  `json:"registered"`
	CompanyType  companies.CompanyType `json:"company_type"`
}

type ModifyCompanyRequest struct {
	Description  string                `json:"description"`
	EmployeesNum uint32                `json:"employees_num"`
	Registered   bool                  `json:"registered"`
	CompanyType  companies.CompanyType `json:"company_type"`
}

// AddCompanyHandler exposes an HTTP endpoint to add new company.
func AddCompanyHandler(client gen.CompaniesClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddCompanyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		cmp, err := client.AddCompany(c.Request.Context(), &gen.AddCompanyRequest{
			Name:         req.Name,
			Description:  req.Description,
			EmployeesNum: req.EmployeesNum,
			Registered:   req.Registered,
			CompanyType:  gen.CompanyType(req.CompanyType),
		})
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		logrus.Infof("company created: %v", cmp)

		c.JSON(http.StatusOK, cmp)
	}
}

// ModifyCompanyHandler exposes an HTTP endpoint to modify an existing company.
func ModifyCompanyHandler(client gen.CompaniesClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ModifyCompanyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		id := c.Param("id")

		cmp, err := client.ModifyCompany(c.Request.Context(), &gen.ModifyCompanyRequest{
			Id:           id,
			Description:  req.Description,
			EmployeesNum: req.EmployeesNum,
			Registered:   req.Registered,
			CompanyType:  gen.CompanyType(req.CompanyType),
		})
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		logrus.Infof("company modified: %v", cmp)

		c.JSON(http.StatusOK, cmp)
	}
}

// DeleteCompanyHandler exposes an HTTP endpoint to delete a company.
func DeleteCompanyHandler(client gen.CompaniesClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := client.DeleteCompany(c.Request.Context(), &gen.DeleteCompanyRequest{
			Id: id,
		})
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		logrus.Infof("company deleted: %v", id)

		c.JSON(http.StatusOK, "")
	}
}

// FindCompanyByIDHandler exposes an HTTP endpoint to search for a single company by an ID.
func FindCompanyByIDHandler(client gen.CompaniesClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		cmp, err := client.FindCompanyByID(c.Request.Context(), &gen.FindCompanyByIDRequest{
			Id: id,
		})
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, cmp)
	}
}
