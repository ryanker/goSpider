function aPost(url, params) {
	return new Promise((resolve, reject) => {
		axios.post(url, params).then(res => {
			const r = res.data;
			if (!r) {
				// console.log('data: ' + r);
				reject('服务器错误，请稍后再试');
			} else if (r.code !== '0') {
				// console.log('message: ' + r.message);
				reject(r.message || '服务器错误，请稍后再试');
			} else {
				resolve(r)
			}
		}).catch(error => {
			console.log(error);
			reject('服务器错误，请稍后再试');
		})
	})
}

function humanSize(n) {
	if (n < 1024) {
		return n + ' B';
	} else if (n < 1024 * 1024) {
		return (n / 1024).toFixed(2) + ' K';
	} else if (n < 1024 * 1024 * 1024) {
		return (n / 1024 / 1024).toFixed(2) + ' M';
	} else if (n < 1024 * 1024 * 1024 * 1024) {
		return (n / 1024 / 1024 / 1024).toFixed(2) + ' G';
	} else if (n < 1024 * 1024 * 1024 * 1024 * 1024) {
		return (n / 1024 / 1024 / 1024 / 1024).toFixed(2) + ' T';
	} else {
		return (n / 1024 / 1024 / 1024 / 1024).toFixed(2) + ' P';
	}
}