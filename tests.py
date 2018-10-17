import time
import mock

import pytest
import requests

SERVER_ADDR = "localhost:8000"
LOGIN_AUTH = ('user@example.com', 'badpassword')


def make_url(path):
    return f"http://{SERVER_ADDR}{path}"


def login():
    user, password = LOGIN_AUTH
    r = requests.post(make_url('/login'), json={'username': user, 'password': password})
    assert r.ok
    return r


def test_can_connect():
    r = requests.get(make_url('/'))
    assert r.ok


def test_can_login_with_correct_credentials():
    user, password = LOGIN_AUTH
    r = requests.post(make_url('/login'), json={'username': user, 'password': password})
    assert r.ok
    assert {
        'success': True,
        'error': None,
        'access_token': mock.ANY,
        'refresh_token': mock.ANY,
    } == r.json()
    assert r.json()['access_token']
    assert r.json()['refresh_token']


def test_cannot_login_with_bad_credentials():
    user, _ = LOGIN_AUTH
    r = requests.post(make_url('/login'), json={'username': user, 'password': f"WRONG"})
    assert 400 == r.status_code
    assert False == r.json()['success']
    assert not r.json()['access_token']
    assert not r.json()['refresh_token']


def test_refresh_works():
    r = login()
    ref_tok = r.json()['refresh_token']
    acc_tok = r.json()['access_token']

    time.sleep(1)  # Ensure time elapses so that expiry changes

    r = requests.post(make_url('/refresh'), json={'refresh_token': ref_tok})
    assert r.ok

    assert True == r.json()['success']
    assert r.json()['access_token']
    assert r.json()['access_token'] != acc_tok


def test_protected_works_with_valid_access_token():
    acc_tok = login().json()['access_token']

    r = requests.get(make_url('/protected'), headers={'authorization': f'bearer {acc_tok}'})
    assert r.ok
    assert 'Hello user@example.com' == r.text


def test_protected_does_not_allow_invalid_token():
    r = requests.get(make_url('/protected'), headers={'authorization': f'bearer BROKENTOKEN'})
    assert 401 == r.status_code
    assert '' == r.text
