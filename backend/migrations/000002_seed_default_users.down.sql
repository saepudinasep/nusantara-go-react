DELETE FROM users
WHERE
    username IN (
        'admin',
        'petugas1',
        'guru1',
        'siswa1'
    );