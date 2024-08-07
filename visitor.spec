# SPDX-FileCopyrightText: The RamenDR authors
# SPDX-License-Identifier: Apache-2.0

Name: visitor
Version: 0.5.0
Release: 1%{?dist}
Summary: Disaster recovery demo

License: Apache-2.0
URL: https://gitlab.com/nirs/visitor

Source: %{name}-%{version}.tar.gz

BuildRequires: golang
BuildRequires: systemd-rpm-macros

Provides: %{name} = %{version}

%global debug_package %{nil}

%description
Disaster recovery demo service

%prep
%autosetup

%build
CGO_ENABLED=0 go build -v -o %{name}

%install
install -D -m 0755 %{name} %{buildroot}%{_bindir}/%{name}
install -D -m 0644 %{name}.service %{buildroot}%{_unitdir}/%{name}.service

%files
%license LICENSE
%{_bindir}/%{name}
%{_unitdir}/%{name}.service

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
* Wed Aug 7 2024 Nir Soffer <nsoffer@redhat.com> - 0.4.0
- Cleanups
* Mon Aug 5 2024 Nir Soffer <nsoffer@redhat.com> - 0.3.0
- Initial package
