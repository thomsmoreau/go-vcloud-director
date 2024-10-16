//go:build functional || openapi || certificate || ALL

/*
 * Copyright 2021 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package govcd

import (
	_ "embed"
	"fmt"

	"github.com/vmware/go-vcloud-director/v3/types/v56"

	. "gopkg.in/check.v1"
)

func (vcd *TestVCD) Test_CertificateInLibrary(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointSSLCertificateLibrary)

	alias := "Test_CertificateInLibrary"

	certificateConfig := &types.CertificateLibraryItem{
		Alias:       alias,
		Certificate: certificate,
	}
	createdCertificate, err := vcd.client.Client.AddCertificateToLibrary(certificateConfig)
	check.Assert(err, IsNil)
	openApiEndpoint, err := getEndpointByVersion(&vcd.client.Client)
	check.Assert(err, IsNil)
	check.Assert(openApiEndpoint, NotNil)
	PrependToCleanupListOpenApi(createdCertificate.CertificateLibrary.Alias, check.TestName(), openApiEndpoint+createdCertificate.CertificateLibrary.Id)

	check.Assert(createdCertificate, NotNil)
	check.Assert(createdCertificate.CertificateLibrary.Id, Not(Equals), "")
	check.Assert(createdCertificate.CertificateLibrary.Alias, Equals, alias)
	check.Assert(createdCertificate.CertificateLibrary.Certificate, Equals, certificate)

	matchesCert, err := createdCertificate.SameAs(certificate)
	check.Assert(err, IsNil)
	check.Assert(matchesCert, Equals, true)

	fetchedCertificate, err := vcd.client.Client.GetCertificateFromLibraryById(createdCertificate.CertificateLibrary.Id)
	check.Assert(err, IsNil)
	check.Assert(fetchedCertificate, NotNil)
	check.Assert(fetchedCertificate.CertificateLibrary.Alias, Equals, alias)
	check.Assert(fetchedCertificate.CertificateLibrary.Certificate, Equals, certificate)

	//test with private key and upload to org context
	adminOrg, err := vcd.client.GetAdminOrgByName(vcd.org.Org.Name)
	check.Assert(err, IsNil)
	check.Assert(adminOrg, NotNil)

	aliasForPrivateKey := "Test_CertificateInLibrary_private_key_test"
	description := "generated by test"

	privateKeyPassphrase := "test"
	certificateWithPrivateKeyConfig := &types.CertificateLibraryItem{
		Alias:                aliasForPrivateKey,
		Certificate:          certificate,
		Description:          description,
		PrivateKey:           privateKey,
		PrivateKeyPassphrase: privateKeyPassphrase,
	}
	createdCertificateWithPrivateKeyConfig, err := adminOrg.AddCertificateToLibrary(certificateWithPrivateKeyConfig)
	check.Assert(err, IsNil)
	openApiEndpoint, err = getEndpointByVersion(&vcd.client.Client)
	check.Assert(err, IsNil)
	check.Assert(openApiEndpoint, NotNil)
	PrependToCleanupListOpenApi(createdCertificateWithPrivateKeyConfig.CertificateLibrary.Alias, check.TestName(),
		openApiEndpoint+createdCertificateWithPrivateKeyConfig.CertificateLibrary.Id)

	check.Assert(createdCertificateWithPrivateKeyConfig, NotNil)
	check.Assert(createdCertificateWithPrivateKeyConfig.CertificateLibrary.Id, Not(Equals), "")
	check.Assert(createdCertificateWithPrivateKeyConfig.CertificateLibrary.Alias, Equals, aliasForPrivateKey)
	check.Assert(createdCertificateWithPrivateKeyConfig.CertificateLibrary.Certificate, Equals, certificate)

	fetchedCertificateWithPrivateKey, err := vcd.client.Client.GetCertificateFromLibraryById(createdCertificateWithPrivateKeyConfig.CertificateLibrary.Id)
	check.Assert(err, IsNil)
	check.Assert(fetchedCertificateWithPrivateKey, NotNil)
	check.Assert(fetchedCertificateWithPrivateKey.CertificateLibrary.Alias, Equals, aliasForPrivateKey)
	check.Assert(fetchedCertificateWithPrivateKey.CertificateLibrary.Certificate, Equals, certificate)

	// check fetching all certificates
	allOrgCertificates, err := adminOrg.GetAllCertificatesFromLibrary(nil)
	check.Assert(err, IsNil)
	check.Assert(allOrgCertificates, NotNil)

	matchingCertificates, err := vcd.client.Client.MatchingCertificatesInLibrary(certificate)
	check.Assert(err, IsNil)
	check.Assert(matchingCertificates, NotNil)

	foundCertificates, err := vcd.client.Client.CountMatchingCertificates(certificate)
	check.Assert(err, IsNil)
	check.Assert(foundCertificates, Equals, len(matchingCertificates))
	check.Assert(foundCertificates, Equals, 1)

	if testVerbose {
		fmt.Printf("(org) how many certificates: %d\n", len(allOrgCertificates))
		for i, oneCertificate := range allOrgCertificates {
			fmt.Printf("%3d %-20s %-53s %s\n", i, oneCertificate.CertificateLibrary.Alias,
				oneCertificate.CertificateLibrary.Id, oneCertificate.CertificateLibrary.Description)
		}
	}
	allExistingCertificates, err := adminOrg.client.GetAllCertificatesFromLibrary(nil)
	check.Assert(err, IsNil)
	check.Assert(allExistingCertificates, NotNil)

	if testVerbose {
		fmt.Printf("(global) how many certificates: %d\n", len(allExistingCertificates))
		for i, oneCertificate := range allExistingCertificates {
			fmt.Printf("%3d %-20s %-53s %s\n", i, oneCertificate.CertificateLibrary.Alias,
				oneCertificate.CertificateLibrary.Id, oneCertificate.CertificateLibrary.Description)
		}
	}

	// check fetching certificate by Name
	foundCertificate, err := vcd.client.Client.GetCertificateFromLibraryByName(alias)
	check.Assert(err, IsNil)
	check.Assert(foundCertificate, NotNil)
	check.Assert(foundCertificate.CertificateLibrary.Alias, Equals, alias)

	foundCertificateWithPrivateKey, err := adminOrg.GetCertificateFromLibraryByName(aliasForPrivateKey)
	check.Assert(err, IsNil)
	check.Assert(foundCertificateWithPrivateKey, NotNil)
	check.Assert(foundCertificateWithPrivateKey.CertificateLibrary.Alias, Equals, aliasForPrivateKey)

	// check update
	newAlias := "newAlias"
	newDescription := "newDescription"
	foundCertificateWithPrivateKey.CertificateLibrary.Alias = newAlias
	foundCertificateWithPrivateKey.CertificateLibrary.Description = newDescription
	updateCertificateWithPrivateKey, err := foundCertificateWithPrivateKey.Update()
	check.Assert(err, IsNil)
	check.Assert(updateCertificateWithPrivateKey, NotNil)
	check.Assert(updateCertificateWithPrivateKey.CertificateLibrary.Alias, Equals, newAlias)
	check.Assert(updateCertificateWithPrivateKey.CertificateLibrary.Description, Equals, newDescription)
	check.Assert(updateCertificateWithPrivateKey.CertificateLibrary.Id, Not(Equals), "")
	check.Assert(updateCertificateWithPrivateKey.CertificateLibrary.Certificate, Equals, certificate)
	check.Assert(updateCertificateWithPrivateKey.CertificateLibrary.PrivateKey, NotNil)           // isn't returned
	check.Assert(updateCertificateWithPrivateKey.CertificateLibrary.PrivateKeyPassphrase, NotNil) // isn't returned

	foundCertificate.CertificateLibrary.Alias = newAlias
	foundCertificate.CertificateLibrary.Description = newDescription
	updateCertificate, err := foundCertificate.Update()
	check.Assert(err, IsNil)
	check.Assert(updateCertificate, NotNil)
	check.Assert(updateCertificate.CertificateLibrary.Alias, Equals, newAlias)
	check.Assert(updateCertificate.CertificateLibrary.Description, Equals, newDescription)
	check.Assert(updateCertificate.CertificateLibrary.Id, Not(Equals), "")
	check.Assert(updateCertificate.CertificateLibrary.Certificate, Equals, certificate)
	check.Assert(updateCertificate.CertificateLibrary.PrivateKey, NotNil)           // isn't returned
	check.Assert(updateCertificate.CertificateLibrary.PrivateKeyPassphrase, NotNil) // isn't returned

	//delete certificate
	err = updateCertificateWithPrivateKey.Delete()
	check.Assert(err, IsNil)
	deletedCertificate, err := vcd.client.Client.GetCertificateFromLibraryById(updateCertificateWithPrivateKey.CertificateLibrary.Id)
	check.Assert(ContainsNotFound(err), Equals, true)
	check.Assert(deletedCertificate, IsNil)

	err = updateCertificate.Delete()
	check.Assert(err, IsNil)
	deletedCertificate, err = adminOrg.client.GetCertificateFromLibraryById(updateCertificate.CertificateLibrary.Id)
	check.Assert(ContainsNotFound(err), Equals, true)
	check.Assert(deletedCertificate, IsNil)

}

func (vcd *TestVCD) Test_GetCertificateFromLibraryByName_ValidatesSymbolsInName(check *C) {
	if vcd.skipAdminTests {
		check.Skip(fmt.Sprintf(TestRequiresSysAdminPrivileges, check.TestName()))
	}
	skipOpenApiEndpointTest(vcd, check, types.OpenApiPathVersion1_0_0+types.OpenApiEndpointSSLCertificateLibrary)

	// When alias contains commas, semicolons, stars, or plus signs, the encoding may reject by the API when we try to Query it
	// Also, spaces present their own issues
	for _, symbol := range []string{";", ",", "+", " ", "*", ":"} {

		alias := fmt.Sprintf("Test%sCertificate%sIn%sLibrary", symbol, symbol, symbol)

		certificateConfig := &types.CertificateLibraryItem{
			Alias:       alias,
			Certificate: certificate,
		}
		createdCertificate, err := vcd.client.Client.AddCertificateToLibrary(certificateConfig)
		check.Assert(err, IsNil)
		openApiEndpoint, err := getEndpointByVersion(&vcd.client.Client)
		check.Assert(err, IsNil)
		check.Assert(openApiEndpoint, NotNil)
		PrependToCleanupListOpenApi(createdCertificate.CertificateLibrary.Alias, check.TestName(),
			openApiEndpoint+createdCertificate.CertificateLibrary.Id)

		check.Assert(createdCertificate, NotNil)
		check.Assert(createdCertificate.CertificateLibrary.Id, Not(Equals), "")
		check.Assert(createdCertificate.CertificateLibrary.Alias, Equals, alias)
		check.Assert(createdCertificate.CertificateLibrary.Certificate, Equals, certificate)

		foundCertificate, err := vcd.client.Client.GetCertificateFromLibraryByName(alias)
		check.Assert(err, IsNil)
		check.Assert(foundCertificate, NotNil)
		check.Assert(foundCertificate.CertificateLibrary.Alias, Equals, alias)

		err = foundCertificate.Delete()
		check.Assert(err, IsNil)
	}
}
