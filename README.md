# git-cl

is a tool to generate a changelog. It works by passing a list of commit messages either through a file
or through stdin. Those will get parsed and outputted as a changelog.

## Usage

Pass a list of commit messages (just the subject) via stdin. They will be transformed to a Markdownfile.

Usage example:

    git log --pretty=format:"%s" | go run main.go

Print out a changelog of all commits since the last tag on the current branch:

    git log --pretty=format:"%s" ( git tag -l --sort=-creatordate  --format='%(if)%(*objectname)%(then)%(*objectname:short)%(else)%(objectname:short)%(end)' | head -n 1 )...HEAD

Print out a changelog of all commits between the last two tags:

    git log --pretty=format:"%s" ( git tag -l --sort=-creatordate  --format='%(if)%(*objectname)%(then)%(*objectname:short)%(else)%(objectname:short)%(end)' | head -n 1 )...( git tag -l --sort=-creatordate  --format='%(if)%(*objectname)%(then)%(*objectname:short)%(else)%(objectname:short)%(end)' | head -n 2 | tail -n 1 )

Note: I'm a fish user, so for bash you probably have to adjust slightly

### Options

#### version

Pass `--version` to provide a version. Will be appended to the header separated by a space

#### print default configuration

Pass `--print-config` to print a default configuration file to stdout. You may use this to start with a custom one

#### apply configuration from file

Pass `--config-file <path-to-file>` to load config from file

### Grouping

In terms of conventional commits the types are mapped to groups like this:

| Type     | Group        |
|----------|--------------|
| feat     | Feature      |
| fix      | Bugfix       |
| doc      | Other        |
| chore    | Other        |
| refactor | Optimization |
| test     | Optimization |
| style    | Usability    |
| perf     | Optimization |
| other    | Other        |

Header is currently a first level headline named `Changelog`.
Footer is currently the text `generated by git-cl`

### Structure of the markdown file

The markdown file consists of a header, the body and the footer. Content besides data derived from the commits
is currently hardcoded.

    HEADER
    
    BODY
    
    FOOTER

Body iterates over the groups and lists the commits for each one