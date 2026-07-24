# Zimbra — UI Framework, Backend Stack, and JAR/WAR Analysis

How the front-end framework and the Java backend are put together, and how to
inspect the archives. This maps the **technology and dependency manifest** — it
does not decompile or reproduce Zimbra's Java source (use the open-source repos
in [08-source-repos.md](08-source-repos.md) for that).

## Front-end framework: AjaxTk / DWT

The webmail and admin UIs are **not** modern SPA frameworks — they are Zimbra's
own **Ajax Toolkit (AjaxTk)**, a JavaScript widget library from the mid-2000s
whose classes are prefixed `Dwt*`. Key traits observed live:

- Widget classes prefixed `Dwt*` (`DwtShell`, `DwtDialog`, `DwtTree`,
  `DwtTreeItem`, `DwtToolBar`, `DwtButton`, `DwtComposite`).
- App classes prefixed `Zm*` (`ZmAppCtxt`, `ZmMailApp`, `ZmZimbraMail`, and for
  admin `ZaApp*`).
- Absolute-positioned, JS-measured layout (why the admin login uses
  `position:absolute; top:40%; margin-top:-160px` rather than fl<exbox — see the
  login-clone notes).
- Theming via **skins** (CSS + `skin.properties` tokens), not utility classes.

Implication for our clone: we reproduce the *look* with modern Vue + Tailwind
and skin tokens; we do **not** and should not adopt AjaxTk.

## Backend stack (the `service` webapp)

`service/WEB-INF/lib` holds **74 JARs (~27 MB)** — the Java mail/provisioning
server. Notable dependencies (third-party libraries; versions are public facts):

| Area | Libraries |
| --- | --- |
| DI / core | Spring 6.0.8 (`spring-core/context/beans/aop/expression`) |
| Search index | Apache **Lucene 3.5.0** (`lucene-core/analyzers/smartcn`) |
| JSON | Jackson 2.9.2 (`databind/core/annotations`, smile) |
| Web / REST | Jetty 9.4.57, Jersey 1.11, `javax.ws.rs` 2.0 |
| Auth | **jjwt 0.7.0** (JWT), oauth 1.4 |
| GraphQL | `graphql-java` 9.0, spqr 0.9.7 |
| Coordination | ZooKeeper 3.4.5 + zkclient |
| SOAP/WS | Apache CXF, wsdl4j, xmlschema, neethi |
| Logging | slf4j 1.7.36, log4j-slf4j 2.17.1 |
| Cache | ehcache 3.1.2 |
| Zimbra's own | `zm-taglib-10.1.17`, `zimlettaglib`, `zm-ews-stub` |

Takeaways relevant to us:
- Zimbra indexes mail with **Lucene** (server-side search), separate from the
  MySQL metadata. Our stack delegates search to IMAP.
- The mail store, SOAP API, and provisioning logic are all Java in
  `zm-mailbox`; there is no way to "reuse" them without running Zimbra.

## How to open / analyze the archives (.jar / .war)

For a future deep dive, standard JVM tooling (no Zimbra-specific tool needed):

```bash
# list contents of a JAR/WAR (class/package names — functional, not source)
jar tf zm-taglib-10.1.17.jar | less
unzip -l some.war

# read the manifest / dependency versions
unzip -p some.jar META-INF/MANIFEST.MF

# decompile a class ONLY for personal study of an open-source project you run
#   (Zimbra FOSS source is public — prefer reading the repos in 08-source-repos.md)
#   tools: CFR (cfr.jar), Procyon, Fernflower/IntelliJ, jadx
java -jar cfr.jar SomeClass.class
```

> Note on licensing: Zimbra's server code is copyrighted (historically ZPL /
> now the Synacor/Zimbra Network + FOSS split). **Do not paste decompiled or
> repo source into this project.** Document structure and behavior; link to the
> upstream source for the code itself.

## WAR/webapp descriptors

Each webapp is a standard Java web application:
- `WEB-INF/web.xml` — servlet/filter mappings, welcome files.
- `WEB-INF/lib/*.jar` — the app's Java classes + dependencies.
- `WEB-INF/classes` / `.tld` — tag library descriptors (`zm-taglib`).
- JSPs at the webapp root render the login and bootstrap pages that then load
  the AjaxTk bundles from `js/`.

## What this means for go-snappymail

| Zimbra | go-snappymail |
| --- | --- |
| AjaxTk/DWT (JS widgets) | Vue 3 + Tailwind |
| Skins (`skin.properties`) | `[data-skin]` token CSS |
| FreeMarker/JSP server render | Vite SPA + Go/Echo JSON API |
| Java `service` webapp (Spring/CXF/Lucene) | Go single binary (Echo + GORM) |
| SOAP + GraphQL API | REST `/api/v1/admin/*` |
| Lucene server search | IMAP search |

We match Zimbra's **appearance and information architecture**, on a completely
different, lighter stack. Nothing in the Java backend is a reuse candidate.
