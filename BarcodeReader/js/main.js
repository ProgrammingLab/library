/*global Quagga,swal,superagent*/
window.onload = () =>{
	class BarcodeReader{
		constructor(w, h, target, callBack){
			this.callBack = callBack;
			this.w = w;
			this.h = h;
			this.target = target;
			Quagga.onDetected(data => this._detectedHandler(data));
			Quagga.onProcessed(data => this._processedHandler(data));
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
			let codeArray = code.toString().split('').map((d)=>parseInt(d, 10));
			if(this._checkDigit(codeArray)){
				this.stop();
				let restart = await this.callBack(code);
				if(restart)this.start();
			}
		}

		_processedHandler(result){
			let drawingCtx = Quagga.canvas.ctx.overlay,
				drawingCanvas = Quagga.canvas.dom.overlay;

			if (result) {
				if (result.boxes) {
					drawingCtx.clearRect(0, 0, parseInt(drawingCanvas.getAttribute("width")), parseInt(drawingCanvas.getAttribute("height")));
					result.boxes.filter(box => {
						return box !== result.box;
					}).forEach(box => {
						Quagga.ImageDebug.drawPath(box, {x: 0, y: 1}, drawingCtx, {color: "green", lineWidth: 2});
					});
				}

				if (result.box) {
					Quagga.ImageDebug.drawPath(result.box, {x: 0, y: 1}, drawingCtx, {color: "blue", lineWidth: 2});
				}

				if (result.codeResult && result.codeResult.code) {
					Quagga.ImageDebug.drawPath(result.line, {x: 'x', y: 'y'}, drawingCtx, {color: 'red', lineWidth: 3});
				}
			}
		}

		_checkDigit(codeArray){
			if(codeArray.length != 13)return false;
			if(
				codeArray[0] != 9 ||
				codeArray[1] != 7 ||
				codeArray[2] < 8
			){
				return false;
			}
			let sum = 0;
			let digit = 0;
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
	class Dataparser{
		static parseNDLXML(xml){
			let parser = new DOMParser();
			let dom = parser.parseFromString(xml, 'text/xml');
			let result = {};
			result.find = dom.getElementsByTagName('numberOfRecords')[0].textContent > 0;
			if(result.find){
				let recodeData = parser.parseFromString(unescape(dom.getElementsByTagName('recordData')[0].textContent), 'text/xml');
				result.title = recodeData.getElementsByTagName('dc:title')[0].textContent;
				result.creator = recodeData.getElementsByTagName('dc:creator')[0].textContent;
			}
			return result;

		}
	}

	let barcodeReader = new BarcodeReader(window.innerWidth, window.innerHeight, '#viewer', async code =>{
		let res;
		try{
			res = await superagent
				.get('http://iss.ndl.go.jp/api/sru?operation=searchRetrieve')
				.query({
					query:`isbn=${code}`
				})
				.timeout({
					response: 5000
				})
			;
		}catch(e){
			console.error(e);
			await swal({
				title: '国会図書館apiからデータを取得できませんでした',
				type: 'error',
				text: e
			});
			return true;
		}
		let bookData = Dataparser.parseNDLXML(res.text);
		let result = await swal({
			type: 'question',
			title: 'これであってる?(自信なさげ)',
			text: bookData.find ? `${bookData.title} / ${bookData.creator}` : `ISBNCode: ${code}`,
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
					type: 'error',
					text: e
				});
				return true;
			}
		}else{
			return true;
		}
	});

	barcodeReader.start();
};