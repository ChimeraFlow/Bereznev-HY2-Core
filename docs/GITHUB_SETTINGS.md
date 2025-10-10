GitHub Settings / Flow

Branch strategy
	•	main — stable, деплой в prod
	•	develop — integration, деплой в staging
	•	feature/* — новые фичи → PR в develop
	•	fix/* — багфиксы → PR в develop
	•	hotfix/* — критические фиксы → PR в main, затем cherry-pick в develop

⚠️ История: используем Squash merge + linear history.
Поэтому не смешиваем это с --no-ff (противоречит linear). Релиз делаем Release-PR develop → main со Squash.

⸻

Branch protection (рекомендованные правила)

Для main и develop
	•	☑ Require a pull request before merging
	•	☑ Require approvals: 1+
	•	☑ Dismiss stale approvals
	•	☑ Require status checks to pass before merging
	•	Укажите ваши CI чек-раннеры, например:
	•	build (Android/JVM/iOS)
	•	unit-tests
	•	lint
	•	publish-dry-run (опционально для develop)
	•	☑ Require linear history
	•	☑ Require conversation resolution
	•	☑ Restrict who can push to matching branches (запрет прямых пушей)

⸻

Pull Requests
	•	Обязательное: 1+ reviewer, все статус-чеки зелёные
	•	Тип мержа: Squash and merge
	•	Title/commits: Conventional Commits
	•	Авто-ассайн ревью: CODEOWNERS

⸻

Commit Style (Conventional Commits)

Типы: feat|fix|chore|docs|refactor|test|ci|build|perf
Примеры:
feat(core): add obfs mode
fix(android): crash in startTun2Socks
ci(release): enable maven central publish


Environments (GitHub → Settings → Environments)
	•	dev — по желанию: превью для feature/* (review apps)
	•	staging — автодеплой из develop
	•	prod — деплой по релиз-тегу из main (vX.Y.Z)

⸻

Release flow
	1.	Открываем Release PR: develop → main (Squash)
	2.	После мержа в main — bump версии и тег: vX.Y.Z
	3.	CI публикует: Maven Central + GitHub Packages
	4.	Back-merge: main → develop или cherry-pick релизного коммита

Рекомендуется автогенерация CHANGELOG (например, Conventional Changelog) и релизных нотсов в GitHub Releases.

⸻

Git команды (шпаргалка)

Фича
git checkout -b feature/obfs-mode
# работа...
git add -A
git commit -m "feat(core): add obfs mode"
git push -u origin feature/obfs-mode
# PR → develop (squash merge)


Хотфикс
git checkout -b hotfix/crash-tun2socks main
# фикс...
git add -A
git commit -m "fix(android): crash in startTun2Socks"
git push -u origin hotfix/crash-tun2socks
# PR → main (squash)
# после мержа:
git checkout develop
git cherry-pick <merge-commit-sha>   # или диапазон


Релиз (через PR!)
# Открываем PR develop → main (Squash & Merge)
# После мержа в main:
git checkout main
git pull
git tag -a v1.0.0 -m "Bereznev HY2 Core 1.0.0"
git push origin v1.0.0
# CI берет тег и публикует


Labels (быстрый набор)

type:feature, type:bug, type:chore, prio:high, prio:normal, good first issue, help wanted, area:core, area:android, area:ci.

⸻

Security & maintenance
	•	SECURITY.md с контактами для приватных репортов уязвимостей
	•	dependabot.yml для автоматических обновлений зависимостей (Gradle/Actions)
	•	В CI спрятать секреты: OSSRH_USERNAME, OSSRH_TOKEN, GPG_PRIVATE_KEY, GPG_PASSPHRASE, GH_PACKAGES_TOKEN