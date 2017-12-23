/*global Quagga,swal,superagent*/
window.onload = () =>{
	class BarcodeReader{
		constructor(w, h, target, callBack){
			this.callBack = callBack;
			this.w = w;
			this.h = h;
			this.target = target;
			Quagga.onDetected(data => this._detectedHandler(data));
		}
		start(){
			Quagga.init({
				numOfWorkers: 4,
				inputStream : {
					name : "Live",
					type : "LiveStream",
					target: this.target,
					constraints: {
						width: {max: this.w},
						height: {max: this.h},
						aspectRatio: 1
					},
				},
				decoder : {
					readers : ["ean_reader"]
				}
			},  e =>{
				if(e){
					console.error(e);
					return false;
				}
				Quagga.start();
			});
		}
		stop(){
			Quagga.stop();
		}
		async _detectedHandler(data){
			let code = data.codeResult.code;
			console.log(`code : ${code}`);
			let codeArray = code.toString().split('').map((d)=>parseInt(d, 10));
			if(this._checkDigit(codeArray)){
				this.stop();
				let restart = await this.callBack(code);
				if(restart)this.start();
			}
		}
		_checkDigit(codeArray){
			let sum = 0;
			let digit = 0;
			if(codeArray.length != 13)return false;
			for(let i = 0;i < 12;i++){
				if(i % 2 == 0){
					sum += codeArray[i];
				}else{
					sum += codeArray[i] * 3;
				}
			}
			if(sum % 10 != 0){
				digit = 10 - sum % 10;
			}
			if(digit != codeArray[12])return false;
			return true;
		}

	}

	let barcodeReader = new BarcodeReader(window.innerWidth, window.innerHeight, '#container', async code =>{
		let result = await swal({
			type: 'question',
			title: 'これであってる?(自信なさげ)',
			text: code,
			showCancelButton: true,
			cancelButtonColor: '#d33',
			cancelButtonText: '合ってない',
			confirmButtonText: '合ってる！'
		});
		if(result.value){
			try{
				await superagent
					.get('./')
					.query({sibn: code})
				;
				await swal({
					title: 'データが送信されました',
					type: 'success'
				});
				return true;
			}catch(e){
				console.error(e);
				await swal({
					title: 'データを送信できませんでした',
					type: 'error'
				});
				return true;
			}
		}else{
			return true;
		}
	});

	barcodeReader.start();
};