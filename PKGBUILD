# Maintainer: Volcha
pkgname=slyrics-git
pkgver=r1.0000000
pkgrel=1
pkgdesc="Synchronized lyrics in your terminal, including YouTube Music"
arch=('x86_64')
url="https://github.com/Volcha-Z/Slyrics"
license=('MIT')
makedepends=('go' 'git')
provides=('slyrics')
conflicts=('slyrics')
source=("$pkgname::git+https://github.com/Volcha-Z/Slyrics.git")
sha256sums=('SKIP')

pkgver() {
  cd "$pkgname"
  printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)"
}

build() {
  cd "$pkgname"
  go build -ldflags '-w -s' -o slyrics-bin .
}

package() {
  cd "$pkgname"
  install -Dm755 slyrics-bin "$pkgdir/usr/bin/slyrics-bin"
  install -Dm755 scripts/slyrics "$pkgdir/usr/bin/slyrics"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
