#!/usr/bin/env bash
#
# git add ./internal/integrationtest/samples/gen-samples.sh && git commit -m wip && git push; ./internal/integrationtest/samples/gen-samples.sh
#
set -Euo pipefail

cd "$(dirname "$0")" || exit 255

cat <<"EOF" > success_SuccessParen.go
package integrationtest

type (
	// SuccessParen
	//
	// spanddl:   table: CREATE TABLE `SuccessParens`
	// spanddl: options: PRIMARY KEY (`Id`)
	SuccessParen struct {
		// ID is a primary key of SuccessParen
		// ID requires uuid
		ID string `dbtest:"Id" spanddl:"STRING(36) NOT NULL"`
		// Description is a description of SuccessParen
		Description string `dbtest:"Description" spanddl:"STRING(1024) NOT NULL"`
	}

	// SuccessParenTwo
	//
	// spanddl:      table: CREATE TABLE `SuccessParenTwos`
	// spanddl: constraint: FOREIGN KEY (`SuccessParenId`) REFERENCES SuccessParens(`Id`)
	// spanddl:    options: PRIMARY KEY (`Id`)
	SuccessParenTwo struct {
		// ID is a primary key of SuccessParenTwo
		// ID requires uuid
		ID string `dbtest:"Id" spanddl:"STRING(36) NOT NULL"`
		// SuccessParenID is a foreign key to SuccessParen.ID
		SuccessParenID string `dbtest:"SuccessParenId" spanddl:"STRING(36) NOT NULL"`
		// Description is a description of SuccessParenTwo
		Description string `dbtest:"Description" spanddl:"STRING(1024) NOT NULL"`
	}

	// SuccessParenThree
	//
	// spanddl: options: PRIMARY KEY (`Id`)
	SuccessParenThree struct {
		// ID is a primary key of SuccessParenThree
		// ID requires uuid
		ID string `dbtest:"Id" spanddl:"STRING(36) NOT NULL"`
		// Description is a description of SuccessParenThree
		Description string `dbtest:"Description" spanddl:"STRING(1024) NOT NULL"`
	}

	// spanddl:      table: `SuccessParenFour`
	// spanddl: constraint: FOREIGN KEY (`SuccessParenId`) REFERENCES SuccessParens(`Id`)
	// spanddl:    options: PRIMARY KEY (`Id`)
	SuccessParenFour struct {
		// ID is a primary key of SuccessParenFour
		// ID requires uuid
		ID string `dbtest:"Id" spanddl:"STRING(36) NOT NULL"`
		// SuccessParenID is a foreign key to SuccessParen.ID
		SuccessParenID string `dbtest:"SuccessParenId" spanddl:"STRING(36) NOT NULL"`
		// Description is a description of SuccessParenFour
		Description string `dbtest:"Description" spanddl:"STRING(1024) NOT NULL"`
	}
)
EOF

cat <<"EOF" > success_SuccessRoot.go
package integrationtest

/*
SuccessRoot

spanddl:   table: CREATE TABLE success_roots
spanddl: options: PRIMARY KEY (`Id`)
*/
type SuccessRoot struct {
	// ID is a primary key of SuccessRoot
	// ID requires uuid
	ID string `dbtest:"Id" spanddl:"STRING(36) NOT NULL"`
	// Description is a description of SuccessRoot
	Description string `dbtest:"Description" spanddl:"STRING(1024) NOT NULL"`
}

// SuccessRootTwo
//
// spanddl:      table: CREATE TABLE success_root_twos
// spanddl: constraint: FOREIGN KEY (success_root_id) REFERENCES success_roots(`Id`)
// spanddl:    options: PRIMARY KEY (`Id`)
type SuccessRootTwo struct {
	// ID is a primary key of SuccessRootTwo
	// ID requires uuid
	ID string `dbtest:"Id" spanddl:"STRING(36) NOT NULL"`
	// SuccessRootID is a foreign key to SuccessRoot.ID
	SuccessRootID string `dbtest:"success_root_id" spanddl:"STRING(36) NOT NULL"`
	// Description is a description of SuccessRootTwo
	Description string `dbtest:"Description" spanddl:"STRING(1024) NOT NULL"`
}
EOF

cat <<"EOF" > ./success_SuccessUnformattedFile.go
package integrationtest

// SuccessUnformattedFile
			// spanddl:table:success_unformatted_files
type
SuccessUnformattedFile struct{// this comment is ignored
//  ID is a primary key of SuccessUnformattedFile
ID string `dbtest:"Id" spanddl:"STRING(36) NOT NULL"`
}
EOF
