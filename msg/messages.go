package msg

const (
	PACKAGE_UNTAGGED_ERROR = `
You are trying to package an untagged commit which is not allowed by default.
 If you mean to package a new version (X.Y.Z) then you should create a new tag by running
	git tag vX.Y.Z
 If you mean to package the '%s' version then you should move that tag to the last commit by running
 	git tag -d v%s
 Then running
 	git tag v%s
 And then try again.`
	PACKAGE_UNTRACKER_ERROR = `
You are trying to package uncommitted modified files which is not allowed by default.
 If you mean to package a new version then you should commit all modified files first running:
	git add %s
 And
	git commit
 Then creating a tag for the new version with:
	git tag vX.Y.Z
 And then try again.
 If you mean to package the '%s' version then you should commit all modified files (as exemplified above) and then move that tag to the last commit by running
	git tag -d v%s
 Then running
	git tag v%s
  And then try again`
	PACKAGE_COMMITED_ERROR = `
You are trying to package uncommitted modified files which is not allowed by default.
 If you mean to package a new version then you should commit all modified files first running:
	git commit
 Then creating a tag for the new version with:
	git tag vX.Y.Z
 And then try again.
 If you mean to package the '%s' version then you should commit all modified files (as exemplified above) and then move that tag to the last commit by running
  	git tag -d v%s
 Then running
	git tag v%s
 And then try again`
	PACKAGE_EVENTUAL_UNTAGGED_ERROR = `
You are trying to package an untagged commit which is not allowed by default.
 If you mean to package a new version (%s-%s-VERSION) then you should create a new tag by running
	git tag %s-%s-VERSION
 If you mean to package the '%s' version then you should move that tag to the last commit by running
 	git tag -d %s
 Then running
 	git tag %s
 And then try again.`
	PACKAGE_EVENTUAL_UNTRACKER_ERROR = `
You are trying to package uncommitted modified files which is not allowed by default.
 If you mean to package a new version then you should commit all modified files first running:
	git add %s
 And
	git commit
 Then creating a tag for the new version with:
	git tag %s-%s-VERSION
 And then try again.
 If you mean to package the '%s' version then you should commit all modified files (as exemplified above) and then move that tag to the last commit by running
	git tag -d %s
 Then running
	git tag %s
  And then try again`
	PACKAGE_EVENTUAL_COMMITED_ERROR = `
You are trying to package uncommitted modified files which is not allowed by default.
 If you mean to package a new version then you should commit all modified files first running:
	git commit
 Then creating a tag for the new version with:
	git tag TRAKER_ID-ISSUE_ID-VERSION
 And then try again.
 If you mean to package the '%s' version then you should commit all modified files (as exemplified above) and then move that tag to the last commit by running
  	git tag -d %s
 Then running
	git tag %s
 And then try again`
	PACKAGE_EVENTUAL_INVALID_TRACKER = `
Issue tracker id '%s' in tag '%s' do not match with project issue tracker id '%s'
 You must create a new tag named '%s' running 
	git tag %s %s
 Then remove the invalid tag running 
	git tag -d %s
 And then try again`
	PACKAGE_EVENTUAL_INVALID_ISSUE = `
Issue id '%s' in tag '%s' do not match with project issue id '%s'
You must create a new tag named '%s' running 
   git tag %s %s
Then remove the invalid tag running 
   git tag -d %s
And then try again`
	PACKAGE_EVENTUAL_FILENAME_WARN = `Avoid use '_' in file name ('%s'). File name recomended: '%s'`
)
