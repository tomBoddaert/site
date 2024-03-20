Output buildServe.gif

Require go
Require npx
Require node

Set Shell "bash"
Set FontSize 30
Set FontFamily "CaskaydiaCove Nerd Font"
Set Height 600
Set Width 1200
Set LineHeight 1.1
Set Theme "Catppuccin Macchiato"
Set Padding 40
Set Margin 0

Hide
  Type "set -e" Enter
  Type "git submodule update" Enter
  Sleep 2s
  Type "tar -zxf ../example/site-example.tar.gz" Enter
  Sleep 1s
  Type "go build .." Enter
  Sleep 2s
  Type "cd site-example && export PATH=..:$PATH && set +e && clear" Enter
Show

Type@100ms "site build serve"
Sleep 500ms Enter

Sleep 5s
Ctrl+C
Sleep 2s
