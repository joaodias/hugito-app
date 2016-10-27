## How to contribute

1. Create an Issue describing the problem you have. If the issue already exists, let people know that you want to work in the issue by saying so :smile:
2. Clone the repository to your computer: `git clone https://github.com/joaodias/hugito-app.git`
3. Create a local branch starting with the issue number and then with some words related to the feature, as an example "1-cool-feature"
4. Push the branch to the remote repository by doing `git push origin 1-cool-feature`. This will add the label `in progress` to your feature branch and will let others know that you are working on the feature.
5. Start to hack in your local branch :boom:. Make sure to run `go fmt` before you open pull request. The guidelines used for commenting are the usual with other go code:
  - Just comment unexported methods.
  - You have to have a good justification to comment on an unexported mehtod. Don't forget that code should speak by itself.
  - Be smart in the names you choose.
  - Be smart when using blank lines. Just when needed.
6. Be sure to write test for your code. It should be tested. And all the tests should pass before you open the pull request. Just run the `go test` command and make sure all the specs run. Hugito uses the framework ginkgo (BDD framework) and also some table testing for some data driven methods. Make sure you are/get familiar with those before you start.
7. Commit your work as you go. Check the commit structure in the `contributing.json` file.
8. Push your branch and open a pull request. The pull request should close an issue (most of the times it should).
9. Get ready for review. Someone should review your work and you should discuss it to make sure we are on the good path.
10. After a discussion you get approved to merge and Travis CI and GitMagic say that you can merge.
11. Merge it! :ship:

Some contribution rules are being forced by GitMagic. Don't be afraid to make a mistake :satisfied:. It is OK. And GitMagic will tell you what is wrong :cop:. The rules are defined in `contributing.json` and you can check the meaning of each rule here: https://gitmagic.io/rules/

