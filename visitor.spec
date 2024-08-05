Name: visitor
Version: 0.3.0
Release: 1
Summary: Disaster recovery demo

License: Apache-2.0
URL: https://gitlab.com/nirs/visitor

Source: %{name}-%{version}.tar.gz

BuildRequires: golang
Requires: firewalld-filesystem

%define debug_package %{nil}

%description
Disaster recovery demo service

%prep
%setup -q

%build
CGO_ENABLED=0 go build

%install
install -D %{name} %{buildroot}%{_bindir}/%{name}
install -D -m 0644 %{name}.service %{buildroot}%{_unitdir}/%{name}.service
install -D -m 0644 %{name}.xml %{buildroot}%{_prefix}/lib/firewalld/services/%{name}.xml

%files
%license LICENSE
%{_bindir}/%{name}
%{_unitdir}/%{name}.service
%{_prefix}/lib/firewalld/services/%{name}.xml

%pre
if ! getent passwd %{name} >/dev/null; then
    useradd \
        --system \
        --user-group \
        --shell /usr/sbin/nologin \
        --home-dir /run/%{name} \
        --comment "visitor demo user" \
        %{name}
fi

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%changelog
* Mon Aug 5 2024 Nir Soffer <nsoffer@redhat.com> - 0.3.0
- Initial package
