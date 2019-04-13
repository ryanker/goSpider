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