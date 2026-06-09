/* ==========================================================================
   Umbra Softworks — Calculator (Interactive Installation)
   Extracted from templates/division.html
   ========================================================================== */
document.addEventListener('DOMContentLoaded', function() {
(function() {
  var container = document.querySelector('.section-calc-sw');
  if (!container) return;
  var WA_PHONE = container.dataset.waPhone || '';
  var conceptSelect = document.getElementById('sw-concept');
  var customizationSelect = document.getElementById('sw-customization');
  var scaleSelect = document.getElementById('sw-scale');
  var durationSlider = document.getElementById('sw-event-duration');
  var durationInput = document.getElementById('sw-event-duration-input');
  var durationDisplay = document.getElementById('sw-event-duration-display');
  var hwToggle = document.getElementById('sw-hardware-toggle');
  var hwField = document.getElementById('sw-hardware-field');
  var hwCost = document.getElementById('sw-hardware-cost');
  var aiSelect = document.getElementById('sw-ai');

  var devCostEl = document.getElementById('sw-dev-cost');
  var aiRowEl = document.getElementById('sw-ai-row');
  var aiTotalEl = document.getElementById('sw-ai-total');
  var setupCostEl = document.getElementById('sw-setup-cost');
  var liveCostEl = document.getElementById('sw-live-cost');
  var riskBufferEl = document.getElementById('sw-risk-buffer');
  var staffingCountEl = document.getElementById('sw-staffing-count');
  var staffingTotalEl = document.getElementById('sw-staffing-total');
  var customizationRowEl = document.getElementById('sw-customization-row');
  var customizationTotalEl = document.getElementById('sw-customization-total');
  var scaleRowEl = document.getElementById('sw-scale-row');
  var scaleTotalEl = document.getElementById('sw-scale-total');
  var hwRowEl = document.getElementById('sw-hardware-row');
  var hwTotalEl = document.getElementById('sw-hardware-total');
  var grandTotalEl = document.getElementById('sw-grand-total');
  var ctaEl = document.getElementById('sw-cta');

  var HOURLY_BASELINE = parseInt(container.dataset.hourlyBaseline) || 16650;
  var FLAT_FEE = parseInt(container.dataset.flatFee) || 300000;
  var devDayRate = (8 * HOURLY_BASELINE) + FLAT_FEE;
  var setupDayRate = (12 * HOURLY_BASELINE) + FLAT_FEE;

  function formatIDR(v) {
    return 'Rp ' + Math.round(v).toLocaleString('id-ID');
  }

  function formatDots(vString) {
    var digits = vString.replace(/\D/g, '');
    return digits.replace(/\B(?=(\d{3})+(?!\d))/g, '.');
  }

  function update() {
    var devDays = 5;
    var setupDays = 1;

    var conceptOpt = conceptSelect.options[conceptSelect.selectedIndex];
    devDays = parseInt(conceptOpt.dataset.devDays) || 5;
    setupDays = parseInt(conceptOpt.dataset.setupDays) || 1;

    if (customizationSelect.value === 'premium') {
      devDays = devDays * 1.5;
    }

    var aiOpt = aiSelect.options[aiSelect.selectedIndex];
    var aiDevDays = parseInt(aiOpt.dataset.devDays) || 0;
    var aiSurcharge = parseInt(aiOpt.dataset.surcharge) || 0;

    devDays = devDays + aiDevDays;

    var eventHoursVal = parseInt(durationInput.value) || 1;

    var daysEquiv = (eventHoursVal / 24).toFixed(1);
    if (daysEquiv.endsWith('.0')) {
      daysEquiv = Math.round(eventHoursVal / 24);
    }
    durationDisplay.textContent = eventHoursVal + ' Jam (' + daysEquiv + ' Hari)';

    var scaleOpt = scaleSelect.options[scaleSelect.selectedIndex];
    var riskFactor = parseFloat(scaleOpt.dataset.riskFactor) || 1.10;

    var sourceHardware = hwToggle.checked;
    var rawHwCost = parseInt(hwCost.value.replace(/\./g, '')) || 0;

    if (sourceHardware) {
      hwField.classList.add('calc-sw-visible');
    } else {
      hwField.classList.remove('calc-sw-visible');
    }

    var estimatedStandbyDays = Math.ceil(eventHoursVal / 24) || 1;
    var devCost = devDays * devDayRate;
    var setupCost = setupDays * setupDayRate;
    var liveCost = (eventHoursVal * HOURLY_BASELINE) + (estimatedStandbyDays * FLAT_FEE);

    var customOpt = customizationSelect.options[customizationSelect.selectedIndex];
    var bespokePremium = parseInt(customOpt.dataset.premiumFee) || 0;
    var corporatePremium = scaleSelect.value === 'corporate' ? (parseInt(container.dataset.corporateFee) || 675000) : 0;

    var totalWorkDays = devDays + setupDays + estimatedStandbyDays;
    var extraStaffCount = Math.floor(totalWorkDays * 2);
    var staffRate = parseInt(container.dataset.staffRate) || 250000;
    var staffingSurcharge = extraStaffCount * staffRate;

    var subtotalLabor = devCost + setupCost + liveCost + bespokePremium + corporatePremium + staffingSurcharge;
    var laborTotalWithRisk = subtotalLabor * riskFactor;
    var riskBufferDisplay = laborTotalWithRisk - subtotalLabor;

    var hwClientCost = 0;
    var hwMarkup = parseFloat(container.dataset.hwMarkup) || 1.15;
    if (sourceHardware && rawHwCost > 0) {
      hwClientCost = rawHwCost * hwMarkup;
    }

    var grandTotal = laborTotalWithRisk + hwClientCost + aiSurcharge;

    devCostEl.textContent = formatIDR(devCost);
    setupCostEl.textContent = formatIDR(setupCost);
    liveCostEl.textContent = formatIDR(liveCost);
    riskBufferEl.textContent = formatIDR(riskBufferDisplay);
    staffingCountEl.textContent = extraStaffCount;
    staffingTotalEl.textContent = formatIDR(staffingSurcharge);

    if (bespokePremium > 0) {
      customizationRowEl.style.display = 'flex';
      customizationTotalEl.textContent = formatIDR(bespokePremium);
    } else {
      customizationRowEl.style.display = 'none';
    }

    if (corporatePremium > 0) {
      scaleRowEl.style.display = 'flex';
      scaleTotalEl.textContent = formatIDR(corporatePremium);
    } else {
      scaleRowEl.style.display = 'none';
    }

    if (sourceHardware && rawHwCost > 0) {
      hwRowEl.style.display = 'flex';
      hwTotalEl.textContent = formatIDR(hwClientCost);
    } else {
      hwRowEl.style.display = 'none';
    }

    if (aiSelect.value !== 'none') {
      aiRowEl.style.display = 'flex';
      aiTotalEl.textContent = formatIDR(aiSurcharge);
    } else {
      aiRowEl.style.display = 'none';
    }

    grandTotalEl.textContent = formatIDR(grandTotal);

    // WhatsApp pre-filled message
    var conceptLabel = conceptSelect.options[conceptSelect.selectedIndex].text;
    var customLabel = customizationSelect.options[customizationSelect.selectedIndex].text;
    var scaleLabel = scaleSelect.options[scaleSelect.selectedIndex].text;
    var aiLabel = aiSelect.options[aiSelect.selectedIndex].text;
    var hwLabel = sourceHardware && rawHwCost > 0 ? 'Ya (Rp ' + rawHwCost.toLocaleString('id-ID') + ')' : 'Tidak';

    var swMsgLines = [
      'Halo Umbra, saya ingin mendiskusikan instalasi interaktif:',
      '',
      '• Tipe Interaksi: ' + conceptLabel,
      '• Kompleksitas Visual: ' + customLabel,
      '• Skala Acara: ' + scaleLabel,
      '• Durasi Acara: ' + eventHoursVal + ' Jam',
      '• Fitur AI: ' + aiLabel,
      '• Sewa Hardware: ' + hwLabel,
      '',
      'Total Estimasi: Rp ' + Math.round(grandTotal).toLocaleString('id-ID')
    ];
    ctaEl.href = 'https://wa.me/' + WA_PHONE + '?text=' + encodeURIComponent(swMsgLines.join('\n'));
  }

  hwCost.addEventListener('input', function() {
    var startPos = this.selectionStart;
    var originalLength = this.value.length;
    var formatted = formatDots(this.value);
    this.value = formatted;
    var cursorCorrection = formatted.length - originalLength;
    this.setSelectionRange(startPos + cursorCorrection, startPos + cursorCorrection);
    update();
  });

  hwToggle.addEventListener('change', function() {
    if (!this.checked) {
      hwCost.value = '0';
    }
    update();
  });

  durationSlider.addEventListener('input', function() {
    durationInput.value = this.value;
    update();
  });

  durationInput.addEventListener('input', function() {
    var cleanVal = this.value.replace(/\D/g, '');
    this.value = cleanVal;
    var numericVal = parseInt(cleanVal) || 1;
    if (numericVal > 72) {
      durationSlider.value = 72;
    } else {
      durationSlider.value = numericVal;
    }
    update();
  });

  durationInput.addEventListener('blur', function() {
    var val = parseInt(this.value) || 16;
    this.value = val;
    if (val > 72) {
      durationSlider.value = 72;
    } else {
      durationSlider.value = val;
    }
    update();
  });

  conceptSelect.addEventListener('change', update);
  customizationSelect.addEventListener('change', update);
  scaleSelect.addEventListener('change', update);
  aiSelect.addEventListener('change', update);

  update();
})();
});
