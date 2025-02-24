import http from 'k6/http';
import { check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import crypto from 'k6/crypto';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const VALID_FORMATS = ['mp3', 'm4a', 'wav', 'flac', 'opus'];
const INVALID_FORMAT = 'xyz';

const TEST_FILES = [
    'testdata/0_george_0.wav',
    'testdata/1_george_0.wav',
    'testdata/2_george_0.wav',
    'testdata/3_george_0.wav',
    'testdata/4_george_0.wav',
    'testdata/5_george_0.wav'
];

const AUDIO_FILES = TEST_FILES.map((filePath) => ({
    path: filePath,
    data: open(filePath, 'b'),
    hash: crypto.sha256(open(filePath, 'b'), 'hex')
}));

export let options = {
    vus: 1,
    iterations: 1
};

export default function () {
    let userId = "user-"+randomString(6);
    let phraseId = "phrase-"+randomString(6);
    let format = getRandomFormat();
    
    let selectedAudio = AUDIO_FILES[Math.floor(Math.random() * AUDIO_FILES.length)];

    let uploadRes = uploadAudio(userId, phraseId, selectedAudio);
    check(uploadRes, {
        'Upload status is 200': (res) => res.status === 200,
        'Upload contains filename': (res) => res.json().filename !== undefined,
    });

    if (uploadRes.status === 200) {
        let downloadRes = downloadAudio(userId, phraseId, format);
        check(downloadRes, {
            'Download status is 200': (res) => res.status === 200,
        });

        let invalidFormatRes = downloadAudio(userId, phraseId, INVALID_FORMAT);
        check(invalidFormatRes, {
            'Invalid format returns 400': (res) => res.status === 400,
        });

        let missingFileRes = downloadAudio(randomString(6), randomString(6), format);
        check(missingFileRes, {
            'Missing file returns 404': (res) => res.status === 404,
        });

        let uploadNoFileRes = http.post(`${BASE_URL}/audio/user/${userId}/phrase/${phraseId}`, {});
        check(uploadNoFileRes, {
            'Upload without file returns 400': (res) => res.status === 400,
        });
    }
}

function uploadAudio(userId, phraseId, audio) {
    let url = `${BASE_URL}/audio/user/${userId}/phrase/${phraseId}`;
    let formData = {
        audio: http.file(audio.data, audio.path.split('/').pop(), 'audio/wav'),
    };
    return http.post(url, formData);
}

function downloadAudio(userId, phraseId, format) {
    let url = `${BASE_URL}/audio/user/${userId}/phrase/${phraseId}/${format}`;
    return http.get(url);
}

function getRandomFormat() {
    return VALID_FORMATS[Math.floor(Math.random() * VALID_FORMATS.length)];
}

function verifyHash(downloadedData, originalHash) {
    return crypto.sha256(downloadedData, 'hex') === originalHash;
}
